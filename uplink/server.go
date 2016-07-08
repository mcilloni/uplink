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
	"log"
	"net"
	"strconv"

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

func (u *Uplink) startDispatcher() {
	u.dispatcher = startDispatcher(u.Logger)

	u.db.Listen("new_messages", func(payload ...string) {
		if len(payload) != 1 {
			panic("totally broken payload received by notification listener")
		}

		msgID, err := strconv.ParseInt(payload[0], 10, 64)
		if err != nil {
			panic("totally broken payload received by notification listener")
		}

		message, err := u.getMessage(msgID)
		if err != nil {
			u.printErr(err)
			return
		}

		members, err := u.getMessageReceivers(message)
		if err != nil {
			u.printErr(err)
			return
		}

		for _, member := range members {
			name, err := u.getUsername(member.UID)
			if err != nil {
				u.printErr(err)
				return
			}

			u.dispatcher.Notify(member.ID, &pd.Notification{
				Type:     pd.Notification_MESSAGE,
				UserName: name,
				ConvId:   message.Conversation,
				Body:     message.Body,
			})
		}
	})

	u.db.Listen("new_join_reqs", func(payload ...string) {
		if len(payload) != 1 {
			panic("totally broken payload received by invite_req listener")
		}

		inviteID, err := strconv.ParseInt(payload[0], 10, 64)
		if err != nil {
			panic("totally broken payload received by invite_req listener")
		}

		u.Println(payload[0])

		invite, err := u.getInvite(inviteID)
		if err != nil {
			u.printErr(err)
			return
		}

		senderName, err := u.getUsername(invite.Sender)
		if err != nil {
			u.printErr(err)
			return
		}

		conversation, err := u.getConversation(invite.Conversation)
		if err != nil {
			u.printErr(err)
			return
		}

		u.dispatcher.Notify(invite.Receiver, &pd.Notification{
			Type:     pd.Notification_JOIN_REQ,
			UserName: senderName,
			ConvId:   conversation.ID,
			ConvName: conversation.Name,
		})
	})

	u.db.Listen("new_join_accs", func(payload ...string) {
		if len(payload) != 1 {
			panic("totally broken payload received by invite_acc listener")
		}

		inviteID, err := strconv.ParseInt(payload[0], 10, 64)
		if err != nil {
			panic("totally broken payload received by invite_acc listener")
		}

		invite, err := u.getInvite(inviteID)
		if err != nil {
			u.printErr(err)
			return
		}

		receiverName, err := u.getUsername(invite.Receiver)
		if err != nil {
			u.printErr(err)
			return
		}

		conversation, err := u.getConversation(invite.Conversation)
		if err != nil {
			u.printErr(err)
			return
		}

		memberships, err := u.getMemberships(conversation.ID)
		if err != nil {
			u.printErr(err)
			return
		}

		notification := &pd.Notification{
			Type:     pd.Notification_JOIN_ACC,
			UserName: receiverName,
			ConvId:   conversation.ID,
			ConvName: conversation.Name,
		}

		for _, member := range memberships {
			u.dispatcher.Notify(member.ID, notification)
		}
	})

	u.db.Listen("new_friendship_reqs", func(payload ...string) {
		if len(payload) != 1 {
			u.Panicln("totally broken payload received by friend_req listener")
		}

		friendID, err := strconv.ParseInt(payload[0], 10, 64)
		if err != nil {
			u.Panicln("totally broken payload received by friend_req listener")
		}

		u.Printf("NEW FRIENDSHIP_REQ %d\n", friendID)

		friendship, err := u.getFriendship(friendID)
		if err != nil {
			u.printErr(err)
			return
		}

		if friendship.Established {
			u.printStack(fmt.Sprintf("newly requested friendship %d is already established", friendID))
			return
		}

		senderName, err := u.getUsername(friendship.Sender)
		if err != nil {
			u.printErr(err)
			return
		}

		notification := &pd.Notification{
			Type:     pd.Notification_FRIENDSHIP_REQ,
			UserName: senderName,
		}

		u.Printf("SENDING NOTIFICATION %v FROM %d TO USER %d\n", notification, friendship.Sender, friendship.Receiver)

		u.dispatcher.Notify(friendship.Receiver, notification)
	})

	u.db.Listen("new_friendships", func(payload ...string) {
		if len(payload) != 1 {
			u.Panicln("totally broken payload received by friend_acc listener")
		}

		friendID, err := strconv.ParseInt(payload[0], 10, 64)
		if err != nil {
			u.Panicln("totally broken payload received by friend_acc listener")
		}

		friendship, err := u.getFriendship(friendID)
		if err != nil {
			u.printErr(err)
			return
		}

		senderName, err := u.getUsername(friendship.Sender)
		if err != nil {
			u.printErr(err)
			return
		}

		receiverName, err := u.getUsername(friendship.Receiver)
		if err != nil {
			u.printErr(err)
			return
		}

		u.dispatcher.Notify(friendship.Sender, &pd.Notification{
			Type:     pd.Notification_FRIENDSHIP_ACC,
			UserName: receiverName,
		})

		u.dispatcher.Notify(friendship.Receiver, &pd.Notification{
			Type:     pd.Notification_FRIENDSHIP_ACC,
			UserName: senderName,
		})
	})
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
