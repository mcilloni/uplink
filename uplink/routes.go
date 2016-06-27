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

	"golang.org/x/net/context"

	pd "github.com/mcilloni/uplink/protodef"
)

type uplinkRoutes struct {
	u *Uplink
}

func newRoute(u *Uplink) *uplinkRoutes {
	return &uplinkRoutes{u: u}
}

func (r *uplinkRoutes) Exists(ctx context.Context, username *pd.Username) (*pd.ErrCodeResp, error) {
	ecm := new(pd.ErrCodeResp)

	found, err := r.u.existsUser(username.Name)

	if err == nil {
		ecm.Response = &pd.ErrCodeResp_Success{Success: found}
	} else {
		ecm.Response = &pd.ErrCodeResp_ErrCode{ErrCode: err.Code()}
	}

	return ecm, nil
}

func (r *uplinkRoutes) LoginExchange(ctx context.Context, stream pd.Uplink_LoginExchangeServer) error {
	step, err := stream.Recv()

	if err == io.EOF {
		return pd.ErrBrokeProto
	}

	if err != nil {
		return err
	}

	step1, ok := step.LoginSteps.(*pd.LoginReq_Step1)

	if !ok {
		resp := &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: pd.ErrCode_EBROKEPROTO,
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	name := step1.Step1.Name
	user, protoErr := r.u.getUser(name)

	if protoErr != nil {
		resp := &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: protoErr.Code(),
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	resp := &pd.LoginResp{LoginSteps: &pd.LoginResp_Step1{Step1: user.ChToken}}
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
		resp = &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: pd.ErrCode_EBROKEPROTO,
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	if !bytes.Equal(step2.Step2.EncChToken, user.EncChToken) {
		resp = &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: pd.ErrCode_EAUTHFAIL,
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	tok, encTok, err := genTok(user.PublicKey)
	if err != nil {
		return err
	}

	resp = &pd.LoginResp{
		LoginSteps: &pd.LoginResp_Step2{
			Step2: &pd.UserInfo{
				PublicKey:     user.PublicKey,
				EncPrivateKey: user.EncPrivateKey,
				EncChToken:    encTok,
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

	step3, ok := step.LoginSteps.(*pd.LoginReq_Step3)
	if !ok {
		resp = &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: pd.ErrCode_EBROKEPROTO,
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	if !bytes.Equal(step3.Step3.FinalChallenge, tok) {
		resp = &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: pd.ErrCode_EAUTHFAIL,
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	session, protoErr := r.u.newSession(user.ID)
	if err != nil {
		resp = &pd.LoginResp{
			LoginSteps: &pd.LoginResp_ErrCode{
				ErrCode: protoErr.Code(),
			},
		}

		if err = stream.Send(resp); err != nil {
			return err
		}

		return nil
	}

	resp = &pd.LoginResp{
		LoginSteps: &pd.LoginResp_Step3{
			Step3: &pd.SessInfo{
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

func (r *uplinkRoutes) NewUser(ctx context.Context, ureq *pd.NewUserReq) error {

}

func (r *uplinkRoutes) Resume(ctx context.Context, ses *pd.SessInfo) error {

}
