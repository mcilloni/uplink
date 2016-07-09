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

	pd "github.com/mcilloni/uplink/protodef"
)

type sink struct {
	UID  int64
	Sink chan *pd.Notification
}

type msg struct {
	UID          int64
	Notification *pd.Notification
}

type dispatcher struct {
	l *log.Logger

	bins    map[int64][]chan<- *pd.Notification
	started bool

	addSinkChan    chan *sink
	notifyChan     chan *msg
	removeSinkChan chan *sink
}

func startDispatcher(l *log.Logger) *dispatcher {
	d := &dispatcher{
		l:              l,
		bins:           make(map[int64][]chan<- *pd.Notification),
		addSinkChan:    make(chan *sink, 100),
		notifyChan:     make(chan *msg, 1000),
		removeSinkChan: make(chan *sink, 100),
	}

	go d.start()

	return d
}

func (d *dispatcher) addSinkInternal(uid int64, sink chan<- *pd.Notification) {
	bin, ok := d.bins[uid]
	if !ok {
		bin = []chan<- *pd.Notification{}
	}

	d.bins[uid] = append(bin, sink)

	sink <- &pd.Notification{Type: pd.Notification_HANDLER_READY}
}

func (d *dispatcher) notifyInternal(uid int64, notif *pd.Notification) {
	if isReservedID(uid) {
		return
	}

	if bin, ok := d.bins[uid]; ok {
		for _, sink := range bin {
			go func(sink chan<- *pd.Notification) {
				sink <- notif
			}(sink)
		}
	}
}

func (d *dispatcher) removeSinkInternal(uid int64, toRemove chan<- *pd.Notification) {
	defer close(toRemove)

	if bin, ok := d.bins[uid]; ok {
		for i, sink := range bin {
			if sink == toRemove {
				d.bins[uid] = append(bin[:i], bin[i+1:]...)

				return
			}
		}
	}
}

func (d *dispatcher) start() {
	if d.started {
		panic("dispatcher already started")
	}

	d.started = true

	for {
		select {
		case sink := <-d.addSinkChan:
			d.addSinkInternal(sink.UID, sink.Sink)

		case msg := <-d.notifyChan:
			d.notifyInternal(msg.UID, msg.Notification)

		case sink := <-d.removeSinkChan:
			d.removeSinkInternal(sink.UID, sink.Sink)
		}
	}
}

func (d *dispatcher) AddSink(uid int64) chan *pd.Notification {
	newSink := make(chan *pd.Notification)
	d.addSinkChan <- &sink{uid, newSink}

	return newSink
}

func (d *dispatcher) Notify(uid int64, notification *pd.Notification) {
	d.notifyChan <- &msg{uid, notification}
}

func (d *dispatcher) RemoveSink(uid int64, oldSink chan *pd.Notification) {
	d.removeSinkChan <- &sink{uid, oldSink}
}
