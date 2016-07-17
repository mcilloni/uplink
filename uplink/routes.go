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
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	pd "github.com/mcilloni/uplink/protodef"
)

type uplinkRoutes struct {
	u *Uplink
}

func newRoute(u *Uplink) *uplinkRoutes {
	return &uplinkRoutes{u: u}
}

func (r *uplinkRoutes) checkSession(ctx context.Context) (int64, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return -1, pd.ErrNoMetadata
	}

	uidStrSlice, okUID := md["uid"]
	sessidSlice, okSessid := md["sessid"]

	if !okUID || !okSessid || len(uidStrSlice) != 1 || len(sessidSlice) != 1 {
		return -1, pd.ErrBrokeProto
	}

	sessid := sessidSlice[0]
	uid, err := strconv.ParseInt(uidStrSlice[0], 10, 64)
	if err != nil || uid < 1 {
		return -1, pd.ErrBrokeProto
	}

	res, err := r.u.checkSession(sessid, uid)
	if err != nil {
		return -1, err
	}

	if !res {
		return -1, pd.ErrNotAuthenticated
	}

	return uid, nil
}

func (r *uplinkRoutes) AcceptFriendship(ctx context.Context, username *pd.Name) (*pd.BoolResp, error) {
	receiverUID, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	sender, err := r.u.getUser(username.Name)
	if err != nil {
		return nil, err
	}

	err = r.u.acceptFriendship(sender.ID, receiverUID)
	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: true}, nil
}

func (r *uplinkRoutes) AcceptInvite(ctx context.Context, convID *pd.ID) (*pd.BoolResp, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	if err = r.u.acceptInvite(uid, convID.Id); err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: true}, nil
}

func (r *uplinkRoutes) ConversationInfo(ctx context.Context, id *pd.ID) (*pd.Conversation, error) {
	_, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	conv, err := r.u.getConversationWithPeek(id.Id)
	if err != nil {
		return nil, err
	}

	return &pd.Conversation{
		Id:   conv.ID,
		Name: conv.Name,
		LastMessage: &pd.Message{
			Tag:        conv.convMsg.Tag,
			SenderName: conv.convMsg.SenderName,
			Timestamp:  conv.convMsg.Timestamp,
			Body:       conv.convMsg.Body,
		},
	}, nil
}

func (r *uplinkRoutes) Conversations(ctx context.Context, _ *pd.Empty) (*pd.ConversationList, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	convs, err := r.u.getConversationsWithPeek(uid)
	if err != nil {
		return nil, err
	}

	ret := &pd.ConversationList{Convs: make([]*pd.Conversation, len(convs))}

	for i, conv := range convs {
		var msg *pd.Message
		if conv.convMsg != nil {
			msg = &pd.Message{
				Tag:        conv.convMsg.Tag,
				SenderName: conv.convMsg.SenderName,
				Timestamp:  conv.convMsg.Timestamp,
				Body:       conv.convMsg.Body,
			}
		}

		ret.Convs[i] = &pd.Conversation{
			Id:          conv.ID,
			Name:        conv.Name,
			LastMessage: msg,
		}
	}

	return ret, nil
}

func (r *uplinkRoutes) Exists(ctx context.Context, username *pd.Name) (*pd.BoolResp, error) {
	if len(username.Name) == 0 {
		return nil, pd.ErrZeroLenArg
	}

	if isReservedName(username.Name) {
		return nil, pd.ErrReservedUser
	}

	found, err := r.u.existsUser(username.Name)

	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: found}, nil
}

func (r *uplinkRoutes) Friends(ctx context.Context, _ *pd.Empty) (*pd.FriendList, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	frienships, err := r.u.getFriendships(uid)
	if err != nil {
		return nil, err
	}

	return &pd.FriendList{
		Friends: frienships,
	}, nil
}

func (r *uplinkRoutes) Invites(ctx context.Context, _ *pd.Empty) (*pd.InviteList, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	convs, err := r.u.getInvites(uid)
	if err != nil {
		return nil, err
	}

	ret := &pd.InviteList{Invites: make([]*pd.Invite, len(convs))}

	for i, conv := range convs {
		ret.Invites[i] = &pd.Invite{
			Who:      conv.SenderName,
			ConvName: conv.ConvName,
			ConvId:   conv.ConvID,
		}
	}

	return ret, nil
}

func (r *uplinkRoutes) Login(_ context.Context, authInfo *pd.AuthInfo) (*pd.SessInfo, error) {
	name := authInfo.Name
	authpass := authInfo.Pass

	if len(name) == 0 || len(authpass) == 0 {
		return nil, pd.ErrZeroLenArg
	}

	user, err := r.u.loginUser(name, authpass)

	if err != nil {
		return nil, err
	}

	session, err := r.u.newSession(user.ID)
	if err != nil {
		return nil, err
	}

	return &pd.SessInfo{
		Uid:       user.ID,
		SessionId: session.SessionID,
	}, nil
}

func (r *uplinkRoutes) Messages(ctx context.Context, opts *pd.FetchOpts) (*pd.MessageList, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	if err = r.u.isMember(uid, opts.ConvId); err != nil {
		return nil, err
	}

	convMsgs, err := r.u.getMessages(opts.ConvId, opts.LastTag)
	if err != nil {
		return nil, err
	}

	lenList := len(convMsgs)

	msgList := &pd.MessageList{
		ConvId:   opts.ConvId,
		Messages: make([]*pd.Message, lenList),
	}

	for i, convMsg := range convMsgs {
		// reversing order (0 is oldest, lenList - 1 newest)
		msgList.Messages[lenList-i-1] = &pd.Message{
			Tag:        convMsg.Tag,
			SenderName: convMsg.SenderName,
			Timestamp:  convMsg.Timestamp,
			Body:       convMsg.Body,
		}
	}

	return msgList, nil
}

func (r *uplinkRoutes) NewConversation(ctx context.Context, convName *pd.Name) (*pd.ID, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	conv, err := r.u.newConversation(uid, convName.Name)
	if err != nil {
		return nil, err
	}

	return &pd.ID{Id: conv.ID}, nil
}

func (r *uplinkRoutes) NewUser(_ context.Context, ureq *pd.AuthInfo) (*pd.SessInfo, error) {
	name, pass := ureq.Name, ureq.Pass

	if len(name) == 0 || len(pass) == 0 {
		return nil, pd.ErrZeroLenArg
	}

	user, err := r.u.register(name, pass)
	if err != nil {
		return nil, err
	}

	session, err := r.u.newSession(user.ID)
	if err != nil {
		return nil, err
	}

	return &pd.SessInfo{
		Uid:       user.ID,
		SessionId: session.SessionID,
	}, nil
}

func (r *uplinkRoutes) Notifications(_ *pd.Empty, stream pd.Uplink_NotificationsServer) error {
	ctx := stream.Context()
	uid, err := r.checkSession(ctx)
	if err != nil {
		return err
	}

	sink := r.u.dispatcher.AddSink(uid)
	defer r.u.dispatcher.RemoveSink(uid, sink)

	for {
		select {
		case notif := <-sink:
			if err := stream.Send(notif); err != nil {
				if err == io.EOF {
					return nil
				}

				return err
			}

		case <-ctx.Done():
			err := ctx.Err()
			if err == nil || err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func (r *uplinkRoutes) Ping(ctx context.Context, _ *pd.Empty) (*pd.BoolResp, error) {
	_, err := r.checkSession(ctx)

	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: true}, nil
}

func (r *uplinkRoutes) ReceivedRequests(ctx context.Context, _ *pd.Empty) (*pd.FriendList, error) {
	senderUID, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	pending, err := r.u.getReceivedRequests(senderUID)
	if err != nil {
		return nil, err
	}

	return &pd.FriendList{Friends: pending}, nil
}

func (r *uplinkRoutes) RequestFriendship(ctx context.Context, username *pd.Name) (*pd.BoolResp, error) {
	senderUID, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	receiver, err := r.u.getUser(username.Name)
	if err != nil {
		return nil, err
	}

	_, err = r.u.newFriendship(senderUID, receiver.ID)
	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: true}, nil
}

func (r *uplinkRoutes) SendInvite(ctx context.Context, invReq *pd.Invite) (*pd.BoolResp, error) {
	senderUID, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	receiver, err := r.u.getUser(invReq.Who)
	if err != nil {
		return nil, err
	}

	_, err = r.u.invite(senderUID, receiver.ID, invReq.ConvId)
	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: true}, nil
}

func (r *uplinkRoutes) SendMessage(ctx context.Context, req *pd.NewMsgReq) (*pd.NewMsgResp, error) {
	senderUID, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := r.u.newMessage(req.ConvId, senderUID, req.Body, true)
	if err != nil {
		return nil, err
	}

	return &pd.NewMsgResp{
		Tag:       msg.Tag,
		Timestamp: msg.RecvTime.Unix(),
	}, nil
}

func (r *uplinkRoutes) SentRequests(ctx context.Context, _ *pd.Empty) (*pd.FriendList, error) {
	senderUID, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	pending, err := r.u.getSentRequests(senderUID)
	if err != nil {
		return nil, err
	}

	return &pd.FriendList{Friends: pending}, nil
}

func (r *uplinkRoutes) SubmitRegID(ctx context.Context, regID *pd.RegID) (*pd.BoolResp, error) {
	uid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.u.newFCMSubscription(uid, regID.RegId)
	if err != nil {
		return nil, err
	}

	if strings.Contains(err.Error(), "ETOOMANYFCMIDS") {
		return nil, pd.ErrTooManyFCMIDs
	}

	return &pd.BoolResp{Success: true}, nil
}
