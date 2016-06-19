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
	"errors"
	"strings"
)

var (
	errAlreadyInvited = errors.New("user already invited")
	errEmptyConv      = errors.New("empty conversation")
	errNoConv         = errors.New("no such conversation")
	errNoUser         = errors.New("no such user")
	errNotInvited     = errors.New("user not invited to the given conversation")
	errNotMember      = errors.New("user not member of conversation")
	errSelfInvite     = errors.New("can't invite yourself")
)

// getMemberships returns all the Member elements "user" belongs to.
func (u *Uplink) getMemberships(user uint64) (convs []Member, err error) {
	err = u.db.Model(Member{}).Where(&Member{UID: user}).Scan(convs)

	return
}

func (u *Uplink) getMessages(conv uint64, limit, offset int) (msgs []Message, err error) {
	err = u.db.Model(Message{}).Where(&Message{Conversation: conv}).Limit(limit).Offset(offset).Scan(msgs)

	return
}

func (u *Uplink) getUser(name string) (user User, err error) {
	err = u.db.Model(User{}).Where(&User{Name: name}).Scan(&user)

	if err == nil && user.ID == 0 {
		err = errNoUser
	}

	return
}

func (u *Uplink) getUsersOf(conv uint64) (users []User, err error) {
	err = u.db.Model(User{}).Joins("JOIN members ON users.id = members.uid").Where("conversation = ?", conv).Scan(users)

	if err == nil && len(users) == 0 {
		err = errNoConv
	}

	return
}

func (u *Uplink) initConversation(keyHash []byte) (conv Conversation, err error) {
	conv.KeyHash = keyHash
	err = u.db.Create(&conv)

	return
}

func (u *Uplink) invite(receiver, sender, convID uint64, recvEncKey []byte) (invite Invite, err error) {
	invite = Invite{
		Conversation: convID,
		Sender:       sender,
		Receiver:     receiver,
		RecvEncKey:   recvEncKey,
	}

	err = u.db.Model(Invite{}).Create(&invite)

	return
}

func (u *Uplink) newMessage(conv uint64, sender uint64, body []byte) (msg Message, err error) {
	msg = Message{
		Conversation: conv,
		Sender:       sender,
		Body:         body,
	}

	err = u.db.Create(&msg)

	if err != nil && strings.Contains(err.Error(), "NOT_MEMBER") {
		err = errNotMember
	}

	return
}

func (u *Uplink) register(name string, pk, epk, encTk []byte, tk string) (user User, err error) {
	user = User{
		Name:          name,
		PublicKey:     pk,
		EncPrivateKey: epk,
		EncChToken:    encTk,
		ChToken:       tk,
	}

	err = u.db.Create(&user)

	return
}

func (u *Uplink) subscribe(user, convID uint64) (member Member, err error) {
	member = Member{UID: user, Conversation: convID}

	err = u.db.Create(&member)

	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "NOT_INVITED"):
			err = errNotInvited
		case strings.Contains(msg, "NO_SELF_INVITE"):
			err = errSelfInvite
		case strings.Contains(msg, "UNIQUE_INVITE"):
			err = errAlreadyInvited
		default:
		}
	}

	return
}
