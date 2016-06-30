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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

func checkKey(pk []byte) bool {
	pubKeyIface, err := x509.ParsePKIXPublicKey(pk)
	if err != nil {
		return false
	}

	_, ok := pubKeyIface.(*rsa.PublicKey)

	return ok
}

func genTok(pk []byte) (tok, encTok []byte, err error) {
	pubKeyIface, _ := x509.ParsePKIXPublicKey(pk) // already checked on insertion
	rsaPubKey, _ := pubKeyIface.(*rsa.PublicKey)  // already checked, too

	tok = make([]byte, 256)

	if _, err = rand.Read(tok); err != nil {
		return
	}

	encTok, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, tok, []byte("ver_tok"))

	return
}
