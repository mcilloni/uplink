CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


DROP OWNED BY uplink;


CREATE TABLE users (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  name CHARACTER VARYING(60) NOT NULL,
  reg_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  authpass TEXT NOT NULL,
  CONSTRAINT NAME_ALREADY_TAKEN UNIQUE (NAME)
);
ALTER TABLE users OWNER TO uplink;

CREATE TABLE conversations (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  name CHARACTER VARYING(200) NOT NULL DEFAULT '',
  creation_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL
);
ALTER TABLE conversations OWNER TO uplink;

CREATE TABLE members (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  uid BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
  join_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  CONSTRAINT ALREADY_MEMBER UNIQUE (uid,conversation)
);
ALTER TABLE members OWNER TO uplink;

CREATE TABLE messages (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
  sender BIGINT NOT NULL REFERENCES users,
  recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  body TEXT NOT NULL
);
ALTER TABLE messages OWNER TO uplink;

CREATE TABLE invites (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
  sender BIGINT NOT NULL REFERENCES users,
  receiver BIGINT NOT NULL REFERENCES users,
  recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  CONSTRAINT UNIQUE_INVITE UNIQUE (conversation,receiver),
  CONSTRAINT NO_SELF_INVITE CHECK (sender <> receiver)
);
ALTER TABLE invites OWNER TO uplink;

CREATE TABLE friendships (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  sender BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  receiver BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  established BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT UNIQUE_FRIENDSHIP UNIQUE (sender,receiver),
  CONSTRAINT NO_SELF_FRIEND CHECK (sender <> receiver)
);
ALTER TABLE friendships OWNER TO uplink;

CREATE TABLE sessions (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  session_id TEXT DEFAULT ENCODE(DIGEST(GEN_RANDOM_BYTES(256), 'SHA256'), 'HEX') NOT NULL,
  uid BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  acc_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  CONSTRAINT UNIQUE_SESSION UNIQUE (session_id)
);
ALTER TABLE sessions OWNER TO uplink;


CREATE OR REPLACE FUNCTION valid_session(_sessid TEXT, _uid BIGINT) RETURNS BOOLEAN
  LANGUAGE plpgsql
  AS $$
  BEGIN
    PERFORM(SELECT 1 FROM sessions S WHERE S.session_id = _sessid AND S.uid = _uid AND S.acc_time > TIMEZONE('utc'::TEXT, (NOW() - (30 * interval '1 day'))));

    IF FOUND THEN
      UPDATE sessions SET acc_time = TIMEZONE('utc'::TEXT, NOW());

      RETURN TRUE;
    ELSE
      RETURN FALSE;
    END IF;
  END
$$;

CREATE OR REPLACE FUNCTION is_special_user(_uid BIGINT, OUT ret BOOLEAN) RETURNS BOOLEAN
  LANGUAGE plpgsql
  AS $$
  BEGIN
    SELECT _uid = 1 INTO ret;
  END
$$;

CREATE OR REPLACE FUNCTION is_member(_uid BIGINT, _convid BIGINT, OUT ret BOOLEAN) RETURNS BOOLEAN
  LANGUAGE plpgsql
  AS $$
  BEGIN
    SELECT EXISTS(SELECT 1 FROM members WHERE conversation = _convid AND uid = _uid) INTO ret;
  END
$$;

CREATE OR REPLACE FUNCTION are_friends(_uid1 BIGINT, _uid2 BIGINT, out ret BOOLEAN) RETURNS BOOLEAN
  LANGUAGE plpgsql
  AS $$
  BEGIN
    SELECT EXISTS(SELECT 1 FROM friendships WHERE (sender = _uid1 AND receiver = _uid2) OR (sender = _uid2 AND receiver = _uid1)) INTO ret;
  END
$$;


CREATE OR REPLACE FUNCTION check_membership() RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
	BEGIN
		IF NOT IS_MEMBER(NEW.sender) AND NOT IS_SPECIAL_USER(NEW.sender) THEN -- 'uplink' has uid = 1 and can do everything
			RAISE EXCEPTION 'NOT_MEMBER';
		END IF;

		RETURN NEW;
	END
$$;

CREATE OR REPLACE FUNCTION check_before_invite() RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
	BEGIN
		IF NOT IS_MEMBER(NEW.sender, NEW.conversation) THEN
			RAISE EXCEPTION 'NOT_MEMBER';
		END IF;

    IF IS_MEMBER(NEW.receiver, NEW.conversation) THEN
      RAISE EXCEPTION 'ALREADY_MEMBER';
    END IF;

    IF NOT ARE_FRIENDS(NEW.sender, NEW.receiver) THEN
      RAISE EXCEPTION 'NOT_FRIENDS';
    END IF;

		RETURN NEW;
	END
$$;

CREATE OR REPLACE FUNCTION remove_empty_conv() RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
	BEGIN
		IF NOT EXISTS(SELECT 1 FROM members WHERE conversation = OLD.conversation) THEN
			DELETE FROM conversations WHERE id = OLD.conversation;
		END IF;

		RETURN NULL;
	END
$$;

CREATE OR REPLACE FUNCTION check_invite() RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
    DECLARE SENDER_ID BIGINT;
	BEGIN
		WITH RK AS (
			DELETE FROM invites
			WHERE conversation = NEW.conversation AND receiver = NEW.uid
			RETURNING sender
		) SELECT RK.sender INTO SENDER_ID FROM RK;

		IF NOT FOUND THEN
			RAISE EXCEPTION 'NOT_INVITED';
		END IF;

    PERFORM PG_NOTIFY('new_join_accs', CAST(SENDER_ID AS TEXT));

		RETURN NEW;
	END
$$;

CREATE OR REPLACE FUNCTION check_friendship() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
    DECLARE TMP RECORD;
  BEGIN
    IF EXISTS(SELECT 1 FROM friendships WHERE sender = NEW.receiver AND receiver = NEW.sender) THEN
      RAISE EXCEPTION 'ALREADY_FRIENDS';
    END IF;

    RETURN NEW;
  END
$$;

CREATE OR REPLACE FUNCTION hash_password() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
  BEGIN
    IF NEW.id <> 1 THEN -- it's pointless to encrypt the 'uplink' password
      SELECT CRYPT(NEW.authpass, GEN_SALT('bf', 10)) INTO NEW.authpass;
    END IF;

    RETURN NEW;
  END
$$;

CREATE OR REPLACE FUNCTION notify_new_message() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
  BEGIN
    PERFORM PG_NOTIFY('new_messages', CAST(NEW.id AS TEXT));

    RETURN NEW;
  END
$$;

CREATE OR REPLACE FUNCTION notify_new_invite() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
  BEGIN
    PERFORM PG_NOTIFY('new_join_reqs', CAST(NEW.id AS TEXT));

    RETURN NEW;
  END
$$;

CREATE OR REPLACE FUNCTION notify_friendship_request() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
  BEGIN
    IF NOT NEW.established THEN
      PERFORM PG_NOTIFY('new_friendship_reqs', CAST(NEW.id AS TEXT));
    END IF;

    RETURN NEW;
  END
$$;

CREATE OR REPLACE FUNCTION notify_friendship_accepted() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
  BEGIN
    IF (NOT OLD.established) AND NEW.established THEN
      PERFORM PG_NOTIFY('new_friendships', CAST(NEW.id AS TEXT));
    END IF;

    RETURN NEW;
  END
$$;

CREATE TRIGGER before_insert_message BEFORE INSERT ON messages FOR EACH ROW EXECUTE PROCEDURE check_membership();
CREATE TRIGGER after_insert_message AFTER INSERT ON messages FOR EACH ROW EXECUTE PROCEDURE notify_new_message();
CREATE TRIGGER after_delete_member AFTER DELETE ON members FOR EACH ROW EXECUTE PROCEDURE remove_empty_conv();
CREATE TRIGGER before_insert_invite BEFORE INSERT ON invites FOR EACH ROW EXECUTE PROCEDURE check_before_invite();
CREATE TRIGGER after_insert_invite AFTER INSERT ON invites FOR EACH ROW EXECUTE PROCEDURE notify_new_invite();
CREATE TRIGGER before_insert_friendship BEFORE INSERT ON friendships FOR EACH ROW EXECUTE PROCEDURE check_friendship();
CREATE TRIGGER after_insert_friendship AFTER INSERT ON friendships FOR EACH ROW EXECUTE PROCEDURE notify_friendship_request();
CREATE TRIGGER after_update_friendship AFTER UPDATE ON friendships FOR EACH ROW EXECUTE PROCEDURE notify_friendship_accepted();
CREATE TRIGGER before_insert_users BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE hash_password();

INSERT INTO users(name, authpass) VALUES ('uplink', '');
