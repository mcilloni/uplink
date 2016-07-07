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

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/galeone/igor"
	pd "github.com/mcilloni/uplink/protodef"
)

func isReservedID(uid int64) bool {
	return uid == 1
}

func isReservedName(name string) bool {
	return name == "uplink"
}

func (u *Uplink) printErr(e error) {
	u.printStack(e.Error())
}

func (u *Uplink) printStack(msg string) {
	debug.PrintStack()
	u.Println(msg)
}

func (u *Uplink) serverFault(e error) error {
	if e != nil {
		u.printErr(e)
		return pd.ServerFault(e)
	}

	return nil
}

func (u *Uplink) connectDB(connStr string) error {
	if u.db == nil {
		db, err := igor.Connect(connStr)

		u.db = db

		return err
	}

	return nil
}

func (u *Uplink) acceptFriendship(sender, receiver int64) error {
	friendship, err := u.getPendingFriendshipOf(sender, receiver)
	if err != nil {
		return err
	}

	if friendship.ID == 0 {
		return pd.ErrNoRequest
	}

	friendship.Established = true

	return u.serverFault(u.db.Updates(friendship))
}

func (u *Uplink) checkSession(sessid string, uid int64) (res bool, err error) {
	if isReservedID(uid) {
		return false, pd.ErrReservedUser
	}

	err = u.serverFault(u.db.Model(Session{}).Select("valid_session(?,?)", sessid, uid).Scan(&res))

	return
}

func (u *Uplink) existsUser(name string) (foundUser bool, err error) {

	_, err = u.getUser(name)

	foundUser = err == nil

	if err == pd.ErrNoUser {
		err = nil
	}

	return
}

func (u *Uplink) getConversation(convID int64) (conv *Conversation, err error) {
	conv = new(Conversation)

	err = u.serverFault(u.db.Model(Conversation{}).Where(&Conversation{ID: convID}).Scan(conv))

	if err == nil && conv.ID == 0 {
		err = pd.ErrNoConv
	}

	return
}

func (u *Uplink) getFriendship(friendID int64) (friendship *Friendship, err error) {
	friendship = new(Friendship)

	err = u.serverFault(u.db.Model(Friendship{}).Where(&Friendship{ID: friendID, Established: true}).Scan(friendship))

	if err == nil && friendship.ID == 0 {
		err = pd.ErrNoFriendship
	}

	return
}

func (u *Uplink) getPendingFriendshipOf(sender, receiver int64) (friendship *Friendship, err error) {
	friendship = new(Friendship)

	err = u.serverFault(u.db.Model(Friendship{}).Where(&Friendship{
		Sender:      sender,
		Receiver:    receiver,
		Established: false,
	}).Scan(friendship))

	if err == nil && friendship.ID == 0 {
		err = pd.ErrNoFriendship
	}

	return
}

// getFriendships returns all the established Friendships of the given user.
func (u *Uplink) getFriendships(user int64) (friendships []string, err error) {
	err = u.serverFault(u.db.CTE(`WITH user_friends AS (
		(SELECT sender AS friend_id FROM friendships WHERE receiver = ? AND established)
		UNION
		(SELECT receiver AS friend_id FROM friendships WHERE sender = ? AND established)
	)`).Table("user_friends AS uf").Select("name").Joins("JOIN users ON users.id = uf.friend_id").Scan(&friendships))

	return
}

func (u *Uplink) getInvite(inviteID int64) (invite *Invite, err error) {
	invite = new(Invite)

	err = u.serverFault(u.db.Model(Invite{}).Where(&Invite{ID: inviteID}).Scan(invite))

	if err == nil && invite.ID == 0 {
		err = pd.ErrNotInvited
	}

	return
}

// getMemberships returns all the Member elements "user" belongs to.
func (u *Uplink) getMemberships(user int64) (convs []Member, err error) {
	err = u.db.Model(Member{}).Where(&Member{UID: user}).Scan(&convs)

	return
}

func (u *Uplink) getMessage(msgID int64) (msg *Message, err error) {
	msg = new(Message)

	err = u.serverFault(u.db.Model(Message{}).Where(&Message{ID: msgID}).Scan(msg))

	return
}

func (u *Uplink) getMessages(conv int64, limit, offset int) (msgs []Message, err error) {
	err = u.serverFault(u.db.Model(Message{}).Where(&Message{Conversation: conv}).Limit(limit).Offset(offset).Scan(msgs))

	return
}

func (u *Uplink) getMessageReceivers(msg *Message) (receivers []Member, err error) {
	members := new(Member).TableName()
	messages := new(Message).TableName()

	err = u.serverFault(u.db.Model(Member{}).Joins(
		fmt.Sprintf("JOIN %s ON %s.conversation = %s.conversation", messages, messages, members),
	).Where(msg).Scan(&receivers))

	return
}

// getPendingFriendships returns all the still pending Friendships of the given user.
func (u *Uplink) getPendingFriendships(user int64) (friendships []string, err error) {
	friendTable := new(Friendship).TableName()
	err = u.serverFault(u.db.Table(friendTable).Select("name").Joins("JOIN users ON users.id = "+friendTable+".receiver").Where("sender = ?", user).Scan(&friendships))

	return
}

func (u *Uplink) getUser(name string) (user *User, err error) {
	if isReservedName(name) {
		return nil, pd.ErrReservedUser
	}

	user = new(User)

	err = u.serverFault(u.db.Model(User{}).Where(&User{Name: name}).Scan(user))

	if err == nil && user.ID == 0 {
		err = pd.ErrNoUser
	}

	return
}

func (u *Uplink) getUserFromID(UID int64) (user *User, err error) {
	user = new(User)

	err = u.serverFault(u.db.Model(User{}).Where(&User{ID: UID}).Scan(user))

	if err == nil && user.ID == 0 {
		err = pd.ErrNoUser
	}

	return
}

func (u *Uplink) getUsername(UID int64) (name string, err error) {
	err = u.serverFault(u.db.Model(User{}).Select("name").Where(&User{ID: UID}).Scan(&name))

	if err == nil && name == "" {
		err = pd.ErrNoUser
	}

	return
}

func (u *Uplink) getUsersOf(conv int64) (users []User, err error) {
	err = u.serverFault(u.db.Model(User{}).Joins("JOIN members ON users.id = members.uid").Where("conversation = ?", conv).Scan(users))

	if err == nil && len(users) == 0 {
		err = pd.ErrNoConv
	}

	return
}

func (u *Uplink) initConversation() (conv *Conversation, err error) {
	conv = new(Conversation)

	err = u.serverFault(u.db.Create(conv))

	return
}

func (u *Uplink) invite(receiver, sender, convID int64) (invite *Invite, err error) {
	if isReservedID(receiver) || isReservedID(sender) {
		return nil, pd.ErrReservedUser
	}

	invite = &Invite{
		Conversation: convID,
		Sender:       sender,
		Receiver:     receiver,
	}

	err = u.serverFault(u.db.Model(Invite{}).Create(invite))

	return
}

func (u *Uplink) loginUser(name, pass string) (user *User, err error) {
	if isReservedName(name) {
		return nil, pd.ErrReservedUser
	}

	user = new(User)

	err = u.serverFault(u.db.Model(User{}).Where("name = ? AND authpass = CRYPT(?, authpass)", name, pass).Scan(user))

	if err == nil && user.ID == 0 {
		err = pd.ErrAuthFail
	}

	return
}

func (u *Uplink) newFriendship(sender int64, receiver int64) (friendship *Friendship, err error) {
	if isReservedID(sender) || isReservedID(receiver) {
		return nil, pd.ErrReservedUser
	}

	friendship = &Friendship{
		Sender:   sender,
		Receiver: receiver,
	}

	err = u.db.Create(friendship)

	if err != nil && strings.Contains(err.Error(), "ALREADY_FRIENDS") {
		return nil, pd.ErrAlreadyFriends
	}

	return friendship, u.serverFault(err)
}

func (u *Uplink) newMessage(conv int64, sender int64, body string) (msg *Message, err error) {
	msg = &Message{
		Conversation: conv,
		Sender:       sender,
		Body:         body,
	}

	err = u.db.Create(msg)

	if err != nil && strings.Contains(err.Error(), "NOT_MEMBER") {
		return nil, pd.ErrNotMember
	}

	return msg, u.serverFault(err)
}

func (u *Uplink) newSession(uid int64) (session *Session, err error) {
	if isReservedID(uid) {
		return nil, pd.ErrReservedUser
	}

	session = &Session{UID: uid, SessionID: uniuri.NewLen(88)}
	e := u.db.Create(session)

	return session, u.serverFault(e)
}

func (u *Uplink) register(name, pass string) (user *User, err error) {
	if isReservedName(name) {
		return nil, pd.ErrReservedUser
	}

	user = &User{
		Name:     name,
		Authpass: pass,
	}

	e := u.db.Create(user)

	if e != nil && strings.Contains(e.Error(), "NAME_ALREADY_TAKEN") {
		return user, pd.ErrNameAlreadyTaken
	}

	return user, u.serverFault(e)
}

func (u *Uplink) subscribe(user, convID int64) (member *Member, err error) {
	member = &Member{UID: user, Conversation: convID}

	e := u.db.Create(member)

	if e != nil {
		msg := e.Error()
		switch {
		case strings.Contains(msg, "NOT_INVITED"):
			err = pd.ErrNotInvited
		case strings.Contains(msg, "NO_SELF_INVITE"):
			err = pd.ErrSelfInvite
		case strings.Contains(msg, "UNIQUE_INVITE"):
			err = pd.ErrAlreadyInvited
		default:
			err = u.serverFault(e)
		}
	}

	return
}
