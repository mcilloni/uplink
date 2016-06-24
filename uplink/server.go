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

	"github.com/galeone/igor"
)

// Uplink instance structure.
type Uplink struct {
	cfg Config
	db  *igor.Database
	log *log.Logger
}

func (u *Uplink) serve(conn net.Conn) {

}

// Start starts a previously configured Uplink instance.
func (u *Uplink) Start() (err error) {
	if err = u.connectDB(u.cfg.ConnInfo); err != nil {
		return
	}

	if listener, err := net.Listen("tcp", u.cfg.ConnInfo); err == nil {
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				return err
			}

			go u.serve(conn)
		}
	}

	return
}

// New Initializes and returns an instance of Uplink according
// to the given Config.
func New(cfg Config, logger *log.Logger) (*Uplink, error) {
	return &Uplink{cfg: cfg, log: logger}, nil
}
