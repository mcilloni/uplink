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

func (r *uplinkRoutes) checkSession(ctx context.Context) (int64, bool, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return -1, false, pd.ErrNoMetadata
	}

	uidStrSlice, okUID := md["uid"]
	sessidSlice, okSessid := md["sessid"]

	if !okUID || !okSessid || len(uidStrSlice) != 1 || len(sessidSlice) != 1 {
		return -1, false, pd.ErrBrokeProto
	}

	sessid := sessidSlice[0]
	uid, err := strconv.ParseInt(uidStrSlice[0], 10, 64)
	if err != nil || uid < 1 {
		return -1, false, pd.ErrBrokeProto
	}

	res, err := r.u.checkSession(sessid, uid)
	if err != nil {
		return -1, false, err
	}

	return uid, res, nil
}

func (r *uplinkRoutes) AcceptFriendship(ctx context.Context, username *pd.Username) (*pd.BoolResp, error) {
	receiverUID, valid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, pd.ErrNotAuthenticated
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

func (r *uplinkRoutes) Exists(ctx context.Context, username *pd.Username) (*pd.BoolResp, error) {
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
	uid, valid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, pd.ErrNotAuthenticated
	}

	frienships, err := r.u.getFriendships(uid)
	if err != nil {
		return nil, err
	}

	return &pd.FriendList{
		Friends: frienships,
	}, nil
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
	uid, valid, err := r.checkSession(ctx)
	if err != nil {
		return err
	}

	if !valid {
		return pd.ErrNotAuthenticated
	}

	sink := make(chan *pd.Notification)
	defer close(sink)

	r.u.dispatcher.addSink(uid, sink)
	defer r.u.dispatcher.removeSink(uid, sink)

	for {
		if err := stream.Send(<-sink); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func (r *uplinkRoutes) PendingFriendships(ctx context.Context, _ *pd.Empty) (*pd.FriendList, error) {
	senderUID, valid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, pd.ErrNotAuthenticated
	}

	pending, err := r.u.getPendingFriendships(senderUID)
	if err != nil {
		return nil, err
	}

	return &pd.FriendList{Friends: pending}, nil
}

func (r *uplinkRoutes) Ping(ctx context.Context, _ *pd.Empty) (*pd.BoolResp, error) {
	_, valid, err := r.checkSession(ctx)

	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: valid}, nil
}

func (r *uplinkRoutes) RequestFriendship(ctx context.Context, username *pd.Username) (*pd.BoolResp, error) {
	senderUID, valid, err := r.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, pd.ErrNotAuthenticated
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
