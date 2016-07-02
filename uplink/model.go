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

func (u *Uplink) connectDB(connStr string) error {
	if u.db == nil {
		db, err := igor.Connect(connStr)

		u.db = db

		return pd.ServerFault(err)
	}

	return nil
}

func (u *Uplink) checkSession(sessid string, uid int64) (res bool, err error) {
	err = pd.ServerFault(u.db.Select("valid_session(?,?)", sessid, uid).Scan(&res))

	return
}

// getMemberships returns all the Member elements "user" belongs to.
func (u *Uplink) getMemberships(user int64) (convs []Member, _ error) {
	err := u.db.Model(Member{}).Where(&Member{UID: user}).Scan(convs)

	return convs, pd.ServerFault(err)
}

func (u *Uplink) getMessages(conv int64, limit, offset int) (msgs []Message, err error) {
	err = pd.ServerFault(u.db.Model(Message{}).Where(&Message{Conversation: conv}).Limit(limit).Offset(offset).Scan(msgs))

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

func (u *Uplink) getUser(name string) (user User, err error) {
	err = pd.ServerFault(u.db.Model(User{}).Where(&User{Name: name}).Scan(&user))

	if err == nil && user.ID == 0 {
		err = pd.ErrNoUser
	}

	return
}

func (u *Uplink) getUsersOf(conv int64) (users []User, err error) {
	err = pd.ServerFault(u.db.Model(User{}).Joins("JOIN members ON users.id = members.uid").Where("conversation = ?", conv).Scan(users))

	if err == nil && len(users) == 0 {
		err = pd.ErrNoConv
	}

	return
}

func (u *Uplink) initConversation(keyHash []byte) (conv Conversation, err error) {
	conv.KeyHash = keyHash
	err = pd.ServerFault(u.db.Create(&conv))

	return
}

func (u *Uplink) invite(receiver, sender, convID int64, recvEncKey []byte) (invite Invite, err error) {
	invite = Invite{
		Conversation: convID,
		Sender:       sender,
		Receiver:     receiver,
		RecvEncKey:   recvEncKey,
	}

	err = pd.ServerFault(u.db.Model(Invite{}).Create(&invite))

	return
}

func (u *Uplink) loginUser(name, pass string) (user User, err error) {
	err = pd.ServerFault(u.db.Model(User{}).Where("name = ? AND authpass = CRYPT(?, authpass)").Scan(&user))

	if err == nil && user.ID == 0 {
		err = pd.ErrAuthFail
	}

	return
}

func (u *Uplink) newMessage(conv int64, sender int64, body []byte) (msg Message, err error) {
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

func (u *Uplink) newSession(uid int64) (session Session, err error) {
	session = Session{UID: uid}
	e := u.db.Create(&session)

	return session, pd.ServerFault(e)
}

func (u *Uplink) register(name, pass string, pk, epk, keyIv, keySalt []byte) (user User, err error) {
	user = User{
		Name:          name,
		Authpass:      pass,
		PublicKey:     pk,
		EncPrivateKey: epk,
		KeyIv:         keyIv,
		KeySalt:       keySalt,
	}

	e := u.db.Create(&user)

	if e != nil && strings.Contains(e.Error(), "NAME_ALREADY_TAKEN") {
		return user, pd.ErrNameAlreadyTaken
	}

	return user, pd.ServerFault(e)
}

func (u *Uplink) subscribe(user, convID int64) (member Member, err error) {
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
