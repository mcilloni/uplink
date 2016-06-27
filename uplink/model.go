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
	"strings"

	"github.com/galeone/igor"
	pd "github.com/mcilloni/uplink/protodef"
)

func (u *Uplink) connectDB(connStr string) *pd.Error {
	if u.db == nil {
		db, err := igor.Connect(connStr)

		u.db = db

		return pd.ServerFault(err)
	}

	return nil
}

func (u *Uplink) checkSession(sessid string, uid int64) (res bool, err *pd.Error) {
	err = pd.ServerFault(u.db.Select("valid_session(?,?)", sessid, uid).Scan(&res))

	return
}

// getMemberships returns all the Member elements "user" belongs to.
func (u *Uplink) getMemberships(user int64) (convs []Member, _ *pd.Error) {
	err := u.db.Model(Member{}).Where(&Member{UID: user}).Scan(convs)

	return convs, pd.ServerFault(err)
}

func (u *Uplink) getMessages(conv int64, limit, offset int) (msgs []Message, err *pd.Error) {
	err = pd.ServerFault(u.db.Model(Message{}).Where(&Message{Conversation: conv}).Limit(limit).Offset(offset).Scan(msgs))

	return
}

func (u *Uplink) existsUser(name string) (foundUser bool, err *pd.Error) {
	_, err = u.getUser(name)

	foundUser = err == nil

	if err == pd.ErrNoUser {
		err = nil
	}

	return
}

func (u *Uplink) getUser(name string) (user User, err *pd.Error) {
	err = pd.ServerFault(u.db.Model(User{}).Where(&User{Name: name}).Scan(&user))

	if err == nil && user.ID == 0 {
		err = pd.ErrNoUser
	}

	return
}

func (u *Uplink) getUsersOf(conv int64) (users []User, err *pd.Error) {
	err = pd.ServerFault(u.db.Model(User{}).Joins("JOIN members ON users.id = members.uid").Where("conversation = ?", conv).Scan(users))

	if err == nil && len(users) == 0 {
		err = pd.ErrNoConv
	}

	return
}

func (u *Uplink) initConversation(keyHash []byte) (conv Conversation, err *pd.Error) {
	conv.KeyHash = keyHash
	err = pd.ServerFault(u.db.Create(&conv))

	return
}

func (u *Uplink) invite(receiver, sender, convID int64, recvEncKey []byte) (invite Invite, err *pd.Error) {
	invite = Invite{
		Conversation: convID,
		Sender:       sender,
		Receiver:     receiver,
		RecvEncKey:   recvEncKey,
	}

	err = pd.ServerFault(u.db.Model(Invite{}).Create(&invite))

	return
}

func (u *Uplink) newMessage(conv int64, sender int64, body []byte) (msg Message, err *pd.Error) {
	msg = Message{
		Conversation: conv,
		Sender:       sender,
		Body:         body,
	}

	e := u.db.Create(&msg)

	if e != nil && strings.Contains(e.Error(), "NOT_MEMBER") {
		return msg, pd.ErrNotMember
	}

	return msg, pd.ServerFault(e)
}

func (u *Uplink) newSession(uid int64) (session Session, err *pd.Error) {
	session = Session{UID: uid}
	e := u.db.Create(&session)

	return session, pd.ServerFault(e)
}

func (u *Uplink) register(name string, pk, epk, encTk []byte, tk string) (user User, err *pd.Error) {
	user = User{
		Name:          name,
		PublicKey:     pk,
		EncPrivateKey: epk,
		EncChToken:    encTk,
		ChToken:       tk,
	}

	e := u.db.Create(&user)

	if e != nil && strings.Contains(e.Error(), "NAME_ALREADY_TAKEN") {
		return user, pd.ErrNameAlreadyTaken
	}

	return user, pd.ServerFault(e)
}

func (u *Uplink) subscribe(user, convID int64) (member Member, err *pd.Error) {
	member = Member{UID: user, Conversation: convID}

	e := u.db.Create(&member)

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
			err = pd.ServerFault(e)
		}
	}

	return
}
