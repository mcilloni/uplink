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

var confPath = flag.String("config", "", "path to the given configuration")

func init() {
	flag.Parse()
}

func main() {
	logger := log.New(os.Stderr, "uplink:", log.LstdFlags)

	if *confPath == "" {
		logger.Fatalln("error: no config provided. use the config parameter to supply one")
	}

	conf, err := uplink.ReadConfig(*confPath)
	if err != nil {
		logger.Fatalf("error: %v\n", err)
	}

	up, err := uplink.New(conf, logger)

	if err != nil {
		logger.Fatal(err)
	}

	if err = up.Start(); err != nil {
		log.Fatal(err)
	}
}
