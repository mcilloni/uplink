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

package main

import (
	"flag"
	"log"
	"os"

	"github.com/mcilloni/uplink/uplink"
)

func main() {
	flag.Parse()
	logger := log.New(os.Stderr, "uplink:", log.LstdFlags)

	up, err := uplink.New(uplink.Config{
		DBConnInfo: "user=uplink password=linkie dbname=uplink sslmode=disable",
	}, logger)

	if err != nil {
		logger.Fatal(err)
	}

	if err = up.Start(); err != nil {
		log.Fatal(err)
	}
}
