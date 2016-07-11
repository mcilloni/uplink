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

const (
	maxMessagesToGet = 20
	reservedUID      = 1
	reservedName     = "uplink"
)

func isReservedID(uid int64) bool {
	return uid == reservedUID
}

func isReservedName(name string) bool {
	return name == reservedName
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

		//db.Log(u.Logger)

		u.db = db

		return err
	}

	return nil
}

func (u *Uplink) acceptInvite(user, convID int64) error {
	membership := &Member{
		UID:          user,
		Conversation: convID,
	}

	if err := u.db.Create(membership); err != nil {
		switch {
		case strings.Contains(err.Error(), "NOT_INVITED"):
			return pd.ErrNotInvited

		default:
			return u.serverFault(err)
		}
	}

	name, err := u.getUsername(user)
	if err != nil {
		return err
	}

	if _, err = u.newMessage(convID, 1, "JOINED:"+name, false); err != nil {
		return err
	}

	return u.notifyNewMember(membership)
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

	if err = u.serverFault(u.db.Updates(friendship)); err != nil {
		return err
	}

	return u.notifyFriendshipEstablished(friendship)
}

func (u *Uplink) checkSession(sessid string, uid int64) (res bool, err error) {
	if isReservedID(uid) {
		return false, pd.ErrReservedUser
	}

	err = u.serverFault(u.db.Model(Session{}).Select("valid_session(?,?)", sessid, uid).Scan(&res))

	return
}

func (u *Uplink) deleteFCMSubscription(regID string) (err error) {
	err = u.serverFault(u.db.Delete(&FCMSubscription{
		RegID: regID,
	}))

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

func (u *Uplink) getConversations(UID int64) (convs []Conversation, err error) {
	conversations := new(Conversation).TableName()
	members := new(Member).TableName()

	err = u.serverFault(u.db.Model(Conversation{}).Joins(
		fmt.Sprintf("JOIN %s ON %s.id = %s.conversation", members, conversations, members),
	).Where(fmt.Sprintf("%s.uid = ?", members), UID).Scan(&convs))

	return
}

func (u *Uplink) getFCMSubscriptions(uid int64) (regids []string, err error) {
	err = u.serverFault(u.db.Model(&FCMSubscription{}).Select("reg_id").
		Where("uid = ?", uid).Scan(&regids),
	)

	return
}

func (u *Uplink) getFCMSubscriptionsForConv(convID int64) (regids []string, err error) {
	members := new(Member).TableName()
	fcmSubscription := new(FCMSubscription).TableName()

	err = u.serverFault(u.db.Model(&FCMSubscription{}).Joins(
		fmt.Sprintf("JOIN %s ON %s.uid = %s.uid", members, fcmSubscription, members),
	).Select("reg_id").Where("conversation = ?", convID).Scan(&regids))

	return
}

func (u *Uplink) getFCMSubscriptionsForUIDList(uids []int64) (regids []string, err error) {
	err = u.serverFault(u.db.Model(&FCMSubscription{}).Select("reg_id").Where("uid IN (?)", uids).Scan(&regids))

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

// getFriendships returns all the established Friendships of the given user.
func (u *Uplink) getFriendships(user int64) (friendships []string, err error) {
	err = u.serverFault(u.db.CTE(`WITH user_friends AS (
		(SELECT sender AS friend_id FROM friendships WHERE receiver = ? AND established)
		UNION
		(SELECT receiver AS friend_id FROM friendships WHERE sender = ? AND established)
	)`, user, user).Table("user_friends AS uf").Select("name").Joins("JOIN users ON users.id = uf.friend_id").Scan(&friendships))

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

type inviteInfo struct {
	SenderName string
	ConvID     int64
	ConvName   string
}

func (u *Uplink) getInvites(uid int64) (convs []inviteInfo, err error) {
	invites := new(Invite).TableName()
	conversations := new(Conversation).TableName()
	users := new(User).TableName()

	err = u.serverFault(u.db.Table(conversations).Joins(
		fmt.Sprintf("JOIN %s ON %s.id = %s.conversation JOIN %s ON %s.sender = %s.id", invites, conversations, invites, users, invites, users),
	).Select(
		fmt.Sprintf("%s.name AS sendername, %s.id as convid, %s.name as convname", users, conversations, conversations),
	).Where(invites+".receiver = ?", uid).Scan(&convs))

	return
}

func (u *Uplink) getMembership(memberID int64) (memb *Member, err error) {
	memb = new(Member)

	err = u.serverFault(u.db.Model(Member{}).Where(&Member{ID: memberID}).Scan(memb))

	if err == nil && memb.ID == 0 {
		err = pd.ErrNotMember
	}

	return
}

// getMemberships returns all the Member elements of a given Conversation.
func (u *Uplink) getMemberships(convID int64) (convs []Member, err error) {
	err = u.db.Model(Member{}).Where(&Member{Conversation: convID}).Scan(&convs)

	return
}

func (u *Uplink) getMessage(msgID int64) (msg *Message, err error) {
	msg = new(Message)

	err = u.serverFault(u.db.Model(Message{}).Where(&Message{ID: msgID}).Scan(msg))

	return
}

type convMsg struct {
	Tag        int64
	SenderName string
	Timestamp  int64
	Body       string
}

// remember: they are returned in reversed order (0 is newest, len - 1 oldest)!
func (u *Uplink) getMessages(conv int64, lastTag int64) (msgs []convMsg, err error) {
	messages := new(Message).TableName()
	users := new(User).TableName()

	whereClause := "conversation = ?"
	whereList := []interface{}{conv}

	if lastTag > 0 {
		whereClause += " AND tag < ?"
		whereList = append(whereList, lastTag)
	}

	err = u.serverFault(u.db.Table(messages).Joins(
		fmt.Sprintf("JOIN %s ON %s.sender = %s.id", users, messages, users),
	).Select(
		"tag, name AS sendername, CAST(EXTRACT(EPOCH FROM recv_time) AS BIGINT) AS timestamp, body",
	).Where(whereClause, whereList...).Limit(maxMessagesToGet).Order("tag DESC").Scan(&msgs))

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

// getReceivedRequests returns all the still pending Friendship requests the user never responded to.
func (u *Uplink) getReceivedRequests(user int64) (friendships []string, err error) {
	friendTable := new(Friendship).TableName()
	err = u.serverFault(u.db.Table(friendTable).Select("name").Joins("JOIN users ON users.id = "+friendTable+".sender").Where("receiver = ? AND NOT established", user).Scan(&friendships))

	return
}

// getSentRequests returns all the still pending Friendship requests the user has sent.
func (u *Uplink) getSentRequests(user int64) (friendships []string, err error) {
	friendTable := new(Friendship).TableName()
	err = u.serverFault(u.db.Table(friendTable).Select("name").Joins("JOIN users ON users.id = "+friendTable+".receiver").Where("sender = ? AND NOT established", user).Scan(&friendships))

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

func (u *Uplink) newConversation(creatorUID int64, name string) (conv *Conversation, err error) {
	conv = &Conversation{Name: name, Creator: creatorUID}

	err = u.serverFault(u.db.Create(conv))
	if err != nil {
		return
	}

	return conv, nil
}

func (u *Uplink) invite(sender, receiver, convID int64) (*Invite, error) {
	if isReservedID(receiver) || isReservedID(sender) {
		return nil, pd.ErrReservedUser
	}

	if sender == receiver {
		return nil, pd.ErrSelfInvite
	}

	invite := &Invite{
		Conversation: convID,
		Sender:       sender,
		Receiver:     receiver,
	}

	err := u.db.Model(Invite{}).Create(invite)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "UNIQUE_INVITE"):
			return nil, pd.ErrAlreadyInvited

		case strings.Contains(err.Error(), "NOT_MEMBER"):
			return nil, pd.ErrNotMember

		case strings.Contains(err.Error(), "ALREADY_MEMBER"):
			return nil, pd.ErrAlreadyMember

		case strings.Contains(err.Error(), "NOT_FRIENDS"):
			return nil, pd.ErrNotFriends

		default:
			return nil, u.serverFault(err)
		}
	}

	if err = u.notifyNewInvite(invite); err != nil {
		return nil, err
	}

	return invite, nil
}

func (u *Uplink) isMember(UID, convID int64) (err error) {
	membership := Member{}
	err = u.serverFault(u.db.Model(Member{}).Where(&Member{
		UID:          UID,
		Conversation: convID,
	}).Scan(&membership))

	if err == nil && membership.ID == 0 {
		err = pd.ErrNotMember
	}

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

func (u *Uplink) newFCMSubscription(uid int64, regid string) (subscription *FCMSubscription, err error) {
	if isReservedID(uid) {
		return nil, pd.ErrReservedUser
	}

	subscription = &FCMSubscription{
		UID:   uid,
		RegID: regid,
	}

	err = u.serverFault(u.db.Create(subscription))

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

	if err = u.db.Create(friendship); err != nil {
		if strings.Contains(err.Error(), "ALREADY_FRIENDS") {
			return nil, pd.ErrAlreadyFriends
		}

		return nil, u.serverFault(err)
	}

	if err = u.notifyFriendshipRequest(friendship); err != nil {
		return nil, err
	}

	return friendship, nil
}

func (u *Uplink) newMessage(conv int64, sender int64, body string, notify bool) (*Message, error) {
	msg := &Message{
		Conversation: conv,
		Sender:       sender,
		Body:         body,
	}

	if err := u.db.Create(msg); err != nil {
		if strings.Contains(err.Error(), "NOT_MEMBER") {
			return nil, pd.ErrNotMember
		}

		u.Printf("FAILED ON MSG %s FROM %d\n", body, sender)

		return nil, u.serverFault(err)
	}

	if notify {
		if err := u.notifyNewMessage(msg); err != nil {
			return nil, err
		}
	}

	return msg, nil
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

func (u *Uplink) updateFCMSubscription(old, new string) (err error) {
	err = u.serverFault(u.db.Where(&FCMSubscription{RegID: old}).Updates(&FCMSubscription{RegID: new}))

	return
}
