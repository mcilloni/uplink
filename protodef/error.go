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

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ServerFault wraps an error (representing an internal server error) into an
// Error.
func ServerFault(e error) error {
	if e == nil {
		return nil
	}

	return grpc.Errorf(codes.Internal, "ESERVERFAULT: server is facing issues")
}

var (
	// ErrAlreadyInvited means that a given user has already been invited to a given
	// conversation.
	ErrAlreadyInvited = grpc.Errorf(codes.AlreadyExists, "EALREADYINVITED: user already invited")

	// ErrAlreadyFriends means that a new friendship relationship between two users
	// cannot be established because they're friends already.
	ErrAlreadyFriends = grpc.Errorf(codes.AlreadyExists, "EALREADYFRIENDS: users already friends")

	// ErrEmptyConv means that a given conversation is empty.
	ErrEmptyConv = grpc.Errorf(codes.NotFound, "EEMPTYCONV: empty conversation")

	// ErrAlreadyMember means that the user can't be invited to the given conversation because it is already a member of it.
	ErrAlreadyMember = grpc.Errorf(codes.AlreadyExists, "EALREADYMEMBER: user already member of the conversation")

	// ErrNameAlreadyTaken means that the wanted username is already taken.
	ErrNameAlreadyTaken = grpc.Errorf(codes.AlreadyExists, "ENAMEALREADYTAKEN: name already taken")

	// ErrNoConv means that there is no such conversation.
	ErrNoConv = grpc.Errorf(codes.NotFound, "ENOCONV: no such conversation")

	// ErrNoUser means that the requested user doesn't exist.
	ErrNoUser = grpc.Errorf(codes.NotFound, "ENOUSER: no such user")

	// ErrNoFriendship means that the requested friendship doesn't exist.
	ErrNoFriendship = grpc.Errorf(codes.NotFound, "ENOFRIENDSHIP: not friend with given user")

	// ErrNoRequest means that the user never received the requested friendship request.
	ErrNoRequest = grpc.Errorf(codes.NotFound, "ENOREQUEST: no friendship request")

	// ErrNotInvited means that the current user has no invite into the conversation.
	ErrNotInvited = grpc.Errorf(codes.PermissionDenied, "ENOTINVITED: user not invited to the given conversation")

	// ErrNotMember means that the user is not member of a given conversation
	ErrNotMember = grpc.Errorf(codes.PermissionDenied, "ENOTMEMBER: user not member of conversation")

	// ErrNotFriends means that the user is not friend with the given user.
	ErrNotFriends = grpc.Errorf(codes.PermissionDenied, "ENOTFRIENDS: user is not a friend")

	// ErrSelfInvite means that the user has tried to invite itself into a
	// conversation.
	ErrSelfInvite = grpc.Errorf(codes.PermissionDenied, "ESELFINVITE: can't invite yourself")

	// ErrBrokeProto means that the user did not follow the protocol correctly.
	ErrBrokeProto = grpc.Errorf(codes.Aborted, "EBROKEPROTO: protocol failure")

	// ErrAuthFail means that the authentication process somehow failed.
	ErrAuthFail = grpc.Errorf(codes.PermissionDenied, "EAUTHFAIL: authentication failure")

	// ErrNotAuthenticated means that the user is not authenticated.
	ErrNotAuthenticated = grpc.Errorf(codes.Unauthenticated, "ENOTAUTHENTICATED: not autenticated")

	// ErrZeroLenArg means that the caller sent an empty parameter when he shouldn't.
	ErrZeroLenArg = grpc.Errorf(codes.InvalidArgument, "EZEROLENARG: argument has zero length")

	// ErrNoMetadata means that the caller forgot to set metadata.
	ErrNoMetadata = grpc.Errorf(codes.InvalidArgument, "ENOMETADATA: no metadata supplied")

	// ErrPassTooShort means that the password given during registration is too short.
	ErrPassTooShort = grpc.Errorf(codes.InvalidArgument, "EPASSTOOSHORT: password too short")

	// ErrReservedUser means that the user tried to access information about a reserved user.
	ErrReservedUser = grpc.Errorf(codes.PermissionDenied, "ERESERVEDUSER: the requested user is reserved")
)
