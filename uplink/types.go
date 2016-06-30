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
	ID            int64 `igor:"primary_key"`
	Authpass      string
	Name          string
	RegTime       time.Time `sql:"default:(now() at time zone 'utc')"`
	PublicKey     []byte
	EncPrivateKey []byte
}

// Conversation represents a conversation between many Users.
type Conversation struct {
	ID           int64 `igor:"primary_key"`
	KeyHash      []byte
	CreationTime time.Time `sql:"default:(now() at time zone 'utc')"`
}

// Member represents the membership of a given User to a Conversation.
type Member struct {
	ID           int64 `igor:"primary_key"`
	UID          int64
	Conversation int64
	JoinTime     time.Time `sql:"default:(now() at time zone 'utc')"`
	EncKey       []byte
}

// Message represents a message belonging to a Conversation.
type Message struct {
	ID           int64 `igor:"primary_key"`
	Conversation int64
	Sender       int64
	RecvTime     time.Time `sql:"default:(now() at time zone 'utc')"`
	Body         []byte
}

// Invite represents an invite to a given Conversation.
type Invite struct {
	ID           int64 `igor:"primary_key"`
	Conversation int64
	Sender       int64
	Receiver     int64
	RecvEncKey   []byte
	RecvTime     time.Time `sql:"default:(now() at time zone 'utc')"`
}

// Friendship represents a relationship between two users that contacted each
// other, and appear in each friendlist.
type Friendship struct {
	ID            int64 `igor:"primary_key"`
	User1         int64
	User2         int64
	EstablishTime time.Time `sql:"default:(now() at time zone 'utc')"`
}

// Session represents a session.
type Session struct {
	ID        int64  `igor:"primary_key"`
	SessionID string `sql:"default:encode(digest(gen_random_bytes(256),'sha256'),'hex')"`
	UID       int64
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

// TableName returns the name of the table associated with Friendship.
func (Friendship) TableName() string {
	return "friendships"
}

// TableName returns the name of the table associated with Session.
func (Session) TableName() string {
	return "sessions"
}
