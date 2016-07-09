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
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/galeone/igor"
	pd "github.com/mcilloni/uplink/protodef"
)

// Uplink instance structure.
type Uplink struct {
	*log.Logger

	cfg        Config
	db         *igor.Database
	dispatcher *dispatcher
}

func (u *Uplink) notifyNewMessage(message *Message) error {
	members, err := u.getMessageReceivers(message)
	if err != nil {
		return err
	}

	for _, member := range members {
		name, err := u.getUsername(member.UID)
		if err != nil {
			return err
		}

		u.dispatcher.Notify(member.ID, &pd.Notification{
			Type:     pd.Notification_MESSAGE,
			UserName: name,
			ConvId:   message.Conversation,
			Body:     message.Body,
		})
	}

	return nil
}

func (u *Uplink) notifyNewInvite(invite *Invite) error {
	senderName, err := u.getUsername(invite.Sender)
	if err != nil {
		return err
	}

	conversation, err := u.getConversation(invite.Conversation)
	if err != nil {
		return err
	}

	u.dispatcher.Notify(invite.Receiver, &pd.Notification{
		Type:     pd.Notification_JOIN_REQ,
		UserName: senderName,
		ConvId:   conversation.ID,
		ConvName: conversation.Name,
	})

	return nil
}

func (u *Uplink) notifyNewMember(member *Member) error {
	newUserName, err := u.getUsername(member.UID)
	if err != nil {
		return err
	}

	conversation, err := u.getConversation(member.Conversation)
	if err != nil {
		return err
	}

	memberships, err := u.getMemberships(conversation.ID)
	if err != nil {
		return err
	}

	notification := &pd.Notification{
		Type:     pd.Notification_JOIN_ACC,
		UserName: newUserName,
		ConvId:   conversation.ID,
		ConvName: conversation.Name,
	}

	for _, member := range memberships {
		u.dispatcher.Notify(member.UID, notification)
	}

	return nil
}

func (u *Uplink) notifyFriendshipRequest(friendship *Friendship) error {
	if friendship.Established {
		return pd.ErrAlreadyFriends
	}

	senderName, err := u.getUsername(friendship.Sender)
	if err != nil {
		return err
	}

	notification := &pd.Notification{
		Type:     pd.Notification_FRIENDSHIP_REQ,
		UserName: senderName,
	}

	u.dispatcher.Notify(friendship.Receiver, notification)

	return nil
}

func (u *Uplink) notifyFriendshipEstablished(friendship *Friendship) error {
	senderName, err := u.getUsername(friendship.Sender)
	if err != nil {
		return err
	}

	receiverName, err := u.getUsername(friendship.Receiver)
	if err != nil {
		return err
	}

	u.dispatcher.Notify(friendship.Sender, &pd.Notification{
		Type:     pd.Notification_FRIENDSHIP_ACC,
		UserName: receiverName,
	})

	u.dispatcher.Notify(friendship.Receiver, &pd.Notification{
		Type:     pd.Notification_FRIENDSHIP_ACC,
		UserName: senderName,
	})

	return nil
}

func (u *Uplink) startDispatcher() {
	u.dispatcher = startDispatcher(u.Logger)
}

// Start starts a previously configured Uplink instance.
func (u *Uplink) Start() (err error) {
	err = u.connectDB(u.cfg.DBConnInfo)

	if err != nil {
		return
	}

	listener, err := net.Listen("tcp", u.cfg.ConnInfo)
	if err != nil {
		return
	}

	defer listener.Close()

	srv := grpc.NewServer()
	pd.RegisterUplinkServer(srv, newRoute(u))

	u.startDispatcher()

	err = srv.Serve(listener)
	if err != nil {
		u.Fatalln(err)
	}

	return
}

// New Initializes and returns an instance of Uplink according
// to the given Config.
func New(cfg Config, logger *log.Logger) (*Uplink, error) {
	return &Uplink{cfg: cfg, Logger: logger}, nil
}
