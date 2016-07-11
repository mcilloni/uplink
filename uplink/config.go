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
	"gopkg.in/gcfg.v1"
)

// Config represents the configuration that should be provided to uplink.Start
// to fully initialize a Server instance.
type Config struct {
	Listener struct {
		Proto    string
		ConnInfo string
	}
	DB struct {
		ConnString string
	}
}

// ReadConfig reads the configuration for uplink from a file specified by path.
// See provided example for syntax.
func ReadConfig(path string) (cfg *Config, err error) {
	cfg = new(Config)
	err = gcfg.ReadFileInto(cfg, path)

	return
}
