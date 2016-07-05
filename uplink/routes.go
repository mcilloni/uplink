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
	"bytes"
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

func (r *uplinkRoutes) Exists(ctx context.Context, username *pd.Username) (*pd.BoolResp, error) {
	if len(username.Name) == 0 {
		return nil, pd.ErrZeroLenArg
	}

	found, err := r.u.existsUser(username.Name)

	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: found}, nil
}

func (r *uplinkRoutes) LoginExchange(stream pd.Uplink_LoginExchangeServer) error {
	step, err := stream.Recv()

	if err == io.EOF {
		return pd.ErrBrokeProto
	}

	if err != nil {
		return err
	}

	step1, ok := step.LoginSteps.(*pd.LoginReq_Step1)

	if !ok {
		return pd.ErrBrokeProto
	}

	name := step1.Step1.Name
	authpass := step1.Step1.Pass

	if len(name) == 0 || len(authpass) == 0 {
		return pd.ErrZeroLenArg
	}

	user, err := r.u.loginUser(name, authpass)

	if err != nil {
		return err
	}

	tok, encTok, err := genTok(user.PublicKey)
	if err != nil {
		return r.u.serverFault(err)
	}

	resp := &pd.LoginResp{
		LoginSteps: &pd.LoginResp_Step1{
			Step1: &pd.LoginAccepted{
				UserInfo: &pd.UserInfo{
					PublicKey:     user.PublicKey,
					EncPrivateKey: user.EncPrivateKey,
					KeyIv:         user.KeyIv,
					KeySalt:       user.KeySalt,
				},
				Challenge: &pd.Challenge{
					Token: encTok,
				},
			},
		},
	}

	if err = stream.Send(resp); err != nil {
		return err
	}

	step, err = stream.Recv()
	if err == io.EOF {
		return pd.ErrBrokeProto
	}

	if err != nil {
		return err
	}

	step2, ok := step.LoginSteps.(*pd.LoginReq_Step2)
	if !ok {
		return pd.ErrBrokeProto
	}

	recvToken := step2.Step2.Token
	if len(recvToken) == 0 {
		return pd.ErrZeroLenArg
	}

	if !bytes.Equal(recvToken, tok) {
		return pd.ErrAuthFail
	}

	session, err := r.u.newSession(user.ID)
	if err != nil {
		return err
	}

	resp = &pd.LoginResp{
		LoginSteps: &pd.LoginResp_Step2{
			Step2: &pd.SessInfo{
				Uid:       user.ID,
				SessionId: session.SessionID,
			},
		},
	}

	if err = stream.Send(resp); err != nil {
		return err
	}

	return nil
}

func (r *uplinkRoutes) NewUser(_ context.Context, ureq *pd.NewUserReq) (*pd.NewUserResp, error) {
	if !checkKey(ureq.PublicKey) {
		return nil, pd.ErrBrokenKey
	}

	name, pass, pk, epk, iv, salt := ureq.Name, ureq.Pass, ureq.PublicKey, ureq.EncPrivateKey, ureq.KeyIv, ureq.KeySalt

	if len(name) == 0 || len(pass) == 0 || len(pk) == 0 || len(iv) == 0 || len(salt) == 0 {
		return nil, pd.ErrZeroLenArg
	}

	user, err := r.u.register(name, pass, pk, epk, iv, salt)
	if err != nil {
		return nil, err
	}

	session, err := r.u.newSession(user.ID)
	if err != nil {
		return nil, err
	}

	return &pd.NewUserResp{
		SessionInfo: &pd.SessInfo{
			Uid:       user.ID,
			SessionId: session.SessionID,
		},
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

func (r *uplinkRoutes) Ping(ctx context.Context, _ *pd.Empty) (*pd.BoolResp, error) {
	_, valid, err := r.checkSession(ctx)

	if err != nil {
		return nil, err
	}

	return &pd.BoolResp{Success: valid}, nil
}
