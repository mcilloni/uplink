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

import "time"

// User represents an user of the service.
type User struct {
	ID            uint64 `igor:"primary_key"`
	Name          string
	RegTime       time.Time `sql:"default:(now() at time zone 'utc')"`
	PublicKey     []byte
	EncPrivateKey []byte
	KeyHash       []byte
}

// Conversation represents a conversation between two Users.
type Conversation struct {
	ID           uint64 `igor:"primary_key"`
	User1        uint64
	User2        uint64
	Key1         []byte
	Key2         []byte
	CreationTime time.Time `sql:"default:(now() at time zone 'utc')"`
}

// Message represents a message belonging to a Conversation.
type Message struct {
	ID           uint64 `igor:"primary_key"`
	Conversation uint64
	RecvTime     time.Time
	Body         []byte
}

// TableName returns the name of the table associated with User.
func (User) TableName() string {
	return "users"
}

// TableName returns the name of the table associated with User.
func (Conversation) TableName() string {
	return "conversations"
}

// TableName returns the name of the table associated with User.
func (Message) TableName() string {
	return "messages"
}
