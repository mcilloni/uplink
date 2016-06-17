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
	"github.com/galeone/igor"
)

// Uplink instance structure.
type Uplink struct {
	cfg Config
	db  *igor.Database
}

func (u *Uplink) connectDB() (err error) {
	if u.db == nil {
		u.db, err = igor.Connect(u.cfg.DBConnInfo)
	}

	return
}

// Start starts a previously configured Uplink instance
func (u *Uplink) Start() (err error) {
	err = u.connectDB()
	return
}

// New Initializes and returns an instance of Uplink according
// to the given Config.
func New(cfg Config) (*Uplink, error) {
	return &Uplink{cfg: cfg}, nil
}
