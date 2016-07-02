CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


DROP OWNED BY uplink;


CREATE TABLE users (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  name CHARACTER VARYING(60) NOT NULL,
  reg_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  public_key BYTEA NOT NULL,
  enc_private_key BYTEA NOT NULL,
  key_salt BYTEA NOT NULL,
  key_iv BYTEA NOT NULL,
  authpass TEXT NOT NULL,
  CONSTRAINT NAME_ALREADY_TAKEN UNIQUE (NAME)
);
ALTER TABLE users OWNER TO uplink;

CREATE TABLE conversations (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  key_hash BYTEA NOT NULL,
  creation_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL
);
ALTER TABLE conversations OWNER TO uplink;

CREATE TABLE members (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  uid BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
  join_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  enc_key BYTEA NOT NULL,
  CONSTRAINT ALREADY_MEMBER UNIQUE (uid,conversation)
);
ALTER TABLE members OWNER TO uplink;

CREATE TABLE messages (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
  sender BIGINT NOT NULL REFERENCES users,
  recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  body BYTEA NOT NULL,
  body_iv BYTEA NOT NULL
);
ALTER TABLE messages OWNER TO uplink;

CREATE TABLE invites (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
  sender BIGINT NOT NULL REFERENCES users,
  receiver BIGINT NOT NULL REFERENCES users,
  recv_enc_key BYTEA NOT NULL,
  recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
  CONSTRAINT UNIQUE_INVITE UNIQUE (conversation,receiver),
  CONSTRAINT NO_SELF_INVITE CHECK (sender <> receiver)
);
ALTER TABLE invites OWNER TO uplink;

CREATE TABLE friendships (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  user1 BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  user2 BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  established BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT UNIQUE_FRIENDSHIP UNIQUE (user1,user2),
  CONSTRAINT NO_SELF_FRIEND CHECK (user1 <> user2)
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


CREATE OR REPLACE FUNCTION check_membership() RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
	BEGIN
		IF NOT EXISTS(SELECT 1 FROM members WHERE conversation = NEW.conversation AND uid = NEW.sender) THEN
			RAISE EXCEPTION 'NOT_MEMBER';
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
	BEGIN
		WITH RK AS (
			DELETE FROM invites
			WHERE conversation = NEW.conversation AND receiver = NEW.uid
			RETURNING recv_enc_key
		) SELECT RK.recv_enc_key INTO NEW.enc_key FROM RK;

		IF NOT FOUND THEN
			RAISE EXCEPTION 'NOT_INVITED';
		END IF;

		RETURN NEW;
	END
$$;

CREATE OR REPLACE FUNCTION invert_values() RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
  BEGIN
    UPDATE NEW SET  user1 = user2, user2 = user1;

    RETURN NEW;
  END
$$;


CREATE TRIGGER before_insert_message BEFORE INSERT ON messages FOR EACH ROW EXECUTE PROCEDURE check_membership();
CREATE TRIGGER after_delete_member AFTER DELETE ON members FOR EACH ROW EXECUTE PROCEDURE remove_empty_conv();
CREATE TRIGGER before_insert_invite BEFORE INSERT ON invites FOR EACH ROW EXECUTE PROCEDURE check_membership();
CREATE TRIGGER before_insert_friendship BEFORE INSERT ON friendships FOR EACH ROW WHEN (NEW.user1 > NEW.user2) EXECUTE PROCEDURE invert_values();


CREATE OR REPLACE FUNCTION valid_session(sessid TEXT, uid BIGINT) RETURNS BOOLEAN
  LANGUAGE plpgsql
  AS $$
  BEGIN
    SELECT 1 FROM sessions S WHERE S.session_id = sessid AND S.uid = uid AND S.acc_time > TIMEZONE('utc'::TEXT, (NOW() - (30 * interval '1 day')));

    IF FOUND THEN
      UPDATE sessions SET acc_time = TIMEZONE('utc'::TEXT, NOW());

      RETURN TRUE;
    ELSE
      RETURN FALSE;
    END IF;
  END
$$;
