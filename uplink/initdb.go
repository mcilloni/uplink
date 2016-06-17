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
    user1 BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    user2 BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    key1 BYTEA NOT NULL,
    key2 BYTEA NOT NULL,
    creation_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
    UNIQUE (id1,id2),
    CONSTRAINT symmetry CHECK (user1 < user2)
  `},
	{"messages", `
    id BIGSERIAL NOT NULL PRIMARY KEY,
    conversation BIGINT NOT NULL REFERENCES conversations ON DELETE CASCADE,
    recv_time TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc'::TEXT, NOW()) NOT NULL,
    body BYTEA NOT NULL
  `},
}

// InitDB initializes the database with the needed tables.
func (u *Uplink) InitDB() error {
	err := u.connectDB()
	if err != nil {
		return err
	}

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
	}

	return tx.Commit()
}
