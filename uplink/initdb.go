/*
 *  uplink, a simple daemon to implement a simple chat protocol
 *  Copyright (C) Marco Cilloni <marco.cilloni@yahoo.com> 2016
 *
 *  This Source Code Form is subject to the terms of the Mozilla Public
 *  License, v. 2.0. If a copy of the MPL was not distributed with this
 *  file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *  Exhibit B is not attached; this software is compatible with the
 *  licenses expressed under Section 1.12 of the MPL v2.
 *
 */

package uplink

import "fmt"

type table struct {
	Name, Body string
}

var tables = []table{
	{"users", `
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name CHARACTER VARYING(60) NOT NULL UNIQUE,
    reg_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
    public_key BYTEA NOT NULL,
    enc_private_key BYTEA NOT NULL,
    key_hash BYTEA NOT NULL
  `},
	{"conversations", `
    id BIGSERIAL NOT NULL PRIMARY KEY,
		key_hash BYTEA NOT NULL,
    creation_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL
  `},
	{"members", `
		id BIGSERIAL NOT NULL PRIMARY KEY,
		uid BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
		conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
		join_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
		enc_key BYTEA NOT NULL,
		CONSTRAINT ALREADY_MEMBER UNIQUE (uid,conversation)
	`},
	{"messages", `
    id BIGSERIAL NOT NULL PRIMARY KEY,
    conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
		sender BIGINT NOT NULL REFERENCES users,
    recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
    body BYTEA NOT NULL
  `},
	{"invites", `
		id BIGSERIAL NOT NULL PRIMARY KEY,
		conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
		sender BIGINT NOT NULL REFERENCES users,
		receiver BIGINT NOT NULL REFERENCES users,
		recv_enc_key BYTEA NOT NULL,
		recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
		CONSTRAINT UNIQUE_INVITE UNIQUE (conversation,receiver),
		CONSTRAINT NO_SELF_INVITE CHECK (sender <> receiver)
	`},
}

const triggerFnFormat = `
	CREATE OR REPLACE FUNCTION %s() RETURNS TRIGGER
		LANGUAGE plpgsql
		AS $$
		BEGIN
			%s
		END
	$$`

var triggerFunctions = map[string]string{
	"check_membership": `
		IF NOT EXISTS(SELECT 1 FROM members WHERE conversation = NEW.conversation AND uid = NEW.sender) THEN
			RAISE EXCEPTION 'NOT_MEMBER';
		END IF;

		RETURN NEW;
	`,
	"remove_empty_conv": `
		IF NOT EXISTS(SELECT 1 FROM members WHERE conversation = OLD.conversation) THEN
			DELETE FROM conversations WHERE id = OLD.conversation;
		END IF;

		RETURN NULL;
  `,
	"check_invite": `
		WITH RK AS (
			DELETE FROM invites
			WHERE conversation = NEW.conversation AND receiver = NEW.uid
			RETURNING recv_enc_key
		) SELECT RK.recv_enc_key INTO NEW.enc_key FROM RK;

		IF NOT FOUND THEN
			RAISE EXCEPTION 'NOT_INVITED';
		END IF;

		RETURN NEW;
	`,
}

var triggers = []string{
	"CREATE TRIGGER before_insert_message BEFORE INSERT ON messages FOR EACH ROW EXECUTE PROCEDURE check_membership()",
	"CREATE TRIGGER after_delete_member AFTER DELETE ON members FOR EACH ROW EXECUTE PROCEDURE remove_empty_conv()",
	"CREATE TRIGGER before_insert_invite BEFORE INSERT ON invites FOR EACH ROW EXECUTE PROCEDURE check_membership()",
}

// InitDB initializes the database with the needed tables.
func (u *Uplink) InitDB() error {
	u.log.Println("connecting to postgresql...")
	err := u.connectDB()
	if err != nil {
		return err
	}

	u.log.Println("connected to database.")

	tx := u.db.Begin()

	for _, table := range tables {
		err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table.Name))
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Exec(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, table.Body))
		if err != nil {
			tx.Rollback()
			return err
		}

		u.log.Println("created table " + table.Name)
	}

	for name, body := range triggerFunctions {
		if err := tx.Exec(fmt.Sprintf(triggerFnFormat, name, body)); err != nil {
			tx.Rollback()
			return err
		}

		u.log.Println("created trigger function " + name)
	}

	for _, trigger := range triggers {
		if err := tx.Exec(trigger); err != nil {
			tx.Rollback()
			return err
		}

	}

	u.log.Println("created triggers.")

	return tx.Commit()
}
