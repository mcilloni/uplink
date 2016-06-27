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

package protodef

// Error type for the Uplink protocol.
type Error struct {
	code ErrCode
	msg  string
}

// Code returns the error code (as defined by proto)
func (e Error) Code() ErrCode {
	return e.code
}

// Error returns the text error, and implements the error interface.
func (e Error) Error() string {
	return e.msg
}

func err(msg string, code ErrCode) *Error {
	return &Error{msg: msg, code: code}
}

// ServerFault wraps an error (representing an internal server error) into an
// Error.
func ServerFault(e error) *Error {
	if e == nil {
		return nil
	}

	return err(e.Error(), ErrCode_ESERVERFAULT)
}

var (
	// ErrAlreadyInvited means that a given has already been invited to a given
	// conversation.
	ErrAlreadyInvited = err("user already invited", ErrCode_EALREADYINVITED)

	// ErrEmptyConv means that a given conversation is empty.
	ErrEmptyConv = err("empty conversation", ErrCode_EEMPTYCONV)

	// ErrNameAlreadyTaken means that the wanted username is already taken.
	ErrNameAlreadyTaken = err("name already taken", ErrCode_ENAMEALREADYTAKEN)

	// ErrNoConv means that there is no such conversation.
	ErrNoConv = err("no such conversation", ErrCode_ENOCONV)

	// ErrNoUser means that the requested user doesn't exist.
	ErrNoUser = err("no such user", ErrCode_ENOUSER)

	// ErrNotInvited means that the current user has no invite into the conversation.
	ErrNotInvited = err("user not invited to the given conversation", ErrCode_ENOTINVITED)

	// ErrNotMember means that the user is not member of a given conversation
	ErrNotMember = err("user not member of conversation", ErrCode_ENOTMEMBER)

	// ErrSelfInvite means that the user has tried to invite itself into a
	// conversation.
	ErrSelfInvite = err("can't invite yourself", ErrCode_ESELFINVITE)

	// ErrBrokeProto means that the user did not follow the protocol correctly.
	ErrBrokeProto = err("protocol failure", ErrCode_EBROKEPROTO)

	// ErrAuthFail means that the authentication process somehow failed.
	ErrAuthFail = err("authentication failure", ErrCode_EAUTHFAIL)
)
