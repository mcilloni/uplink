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

// Conversation represents a conversation between many Users.
type Conversation struct {
	ID           uint64 `igor:"primary_key"`
	KeyHash      []byte
	CreationTime time.Time `sql:"default:(now() at time zone 'utc')"`
}

// Member represents the membership of a given User to a Conversation.
type Member struct {
	ID           uint64 `igor:"primary_key"`
	UID          uint64
	Conversation uint64
	JoinTime     time.Time `sql:"default:(now() at time zone 'utc')"`
	EncKey       []byte
}

// Message represents a message belonging to a Conversation.
type Message struct {
	ID           uint64 `igor:"primary_key"`
	Conversation uint64
	Sender       uint64
	RecvTime     time.Time `sql:"default:(now() at time zone 'utc')"`
	Body         []byte
}

// Invite represents an invite to a given Conversation.
type Invite struct {
	ID           uint64 `igor:"primary_key"`
	Conversation uint64
	Sender       uint64
	Receiver     uint64
	RecvEncKey   []byte
	RecvTime     time.Time `sql:"default:(now() at time zone 'utc')"`
}

// TableName returns the name of the table associated with User.
func (User) TableName() string {
	return "users"
}

// TableName returns the name of the table associated with Conversation.
func (Conversation) TableName() string {
	return "conversations"
}

// TableName returns the name of the table associated with Member.
func (Member) TableName() string {
	return "members"
}

// TableName returns the name of the table associated with Message.
func (Message) TableName() string {
	return "messages"
}

// TableName returns the name of the table associated with Invite.
func (Invite) TableName() string {
	return "invites"
}
