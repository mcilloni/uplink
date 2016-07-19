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
	"unicode/utf8"

	"github.com/golang/protobuf/jsonpb"

	fcm "github.com/mcilloni/go-fcm"
	pd "github.com/mcilloni/uplink/protodef"
)

const (
	maxBodyRunes  = 500
	maxTentatives = 5
)

type regTentatives struct {
	Regids []string
	Trys   []uint8
}

func truncate(str string) string {
	if utf8.RuneCountInString(str) > maxBodyRunes {
		return string([]rune(str)[0:maxBodyRunes])
	}

	return str
}

func notificationToJSON(n *pd.Notification) (string, error) {
	newN := *n
	newN.Body = truncate(n.Body)
	return (&jsonpb.Marshaler{OrigName: true}).MarshalToString(&newN)
}

func (u *Uplink) sendFCMMessage(regids regTentatives, message string) error {
	resp, err := fcm.SendHttp(u.cfg.FCM.APIKey, fcm.HttpMessage{})

	if err != nil {
		return err
	}

	if resp.Failure == 0 && resp.CanonicalIds == 0 {
		return nil
	}

	retryIDs := []string{}
	retryTs := []uint8{}

	for i, result := range resp.Results {
		curRegID := regids.Regids[i]

		if result.MessageId != "" {
			if result.RegistrationId != "" {
				err := u.updateFCMSubscription(curRegID, result.RegistrationId)
				if err != nil {
					u.Printf("CANNOT UPDATE REGID %s : %v\n", curRegID, err)
				}
			}
		} else {
			switch result.Error {
			case "NotRegistered":
				err := u.deleteFCMSubscription(curRegID)
				if err != nil {
					u.Printf("CANNOT DELETE REGID %s : %v\n", curRegID, err)
				}

			case "Unavailable":
				if curTry := regids.Trys[i]; curTry < maxTentatives {
					retryIDs = append(retryIDs, curRegID)
					retryTs = append(retryTs, curTry+1)
				}
			default:
				u.Printf("ERROR IN GCM RESPONSE: %s\n", result.Error)
			}
		}
	}

	if len(retryIDs) > 0 {
		return u.sendFCMMessage(regTentatives{
			Regids: retryIDs,
			Trys:   retryTs,
		}, message)
	}

	return nil
}

func (u *Uplink) sendFCMBroadcast(b *broadcast) error {
	regids, err := u.getFCMSubscriptionsForUIDList(b.UIDs)
	if err != nil {
		return err
	}

	marsh, err := notificationToJSON(b.Notification)
	if err != nil {
		return err
	}

	return u.sendFCMMessage(regTentatives{
		Regids: regids,
		Trys:   make([]uint8, len(regids)),
	}, marsh)
}

func (u *Uplink) registerFCMHandler() {
	sink := make(chan *broadcast, 100)
	u.dispatcher.addPushChan(sink)

	go func() {
		for {
			var err error
			bcast := <-sink

			switch bcast.Notification.Type {
			case pd.Notification_MESSAGE:
				fallthrough
			case pd.Notification_FRIENDSHIP_REQ:
				fallthrough
			case pd.Notification_FRIENDSHIP_ACC:
				fallthrough
			case pd.Notification_JOIN_REQ:
				fallthrough
			case pd.Notification_JOIN_ACC:
				fallthrough
			case pd.Notification_HANDLER_READY:
				err = u.sendFCMBroadcast(bcast)

			default:
				panic("INVALID NOTIFICATION TYPE")
			}

			if err != nil {
				u.Printf("ERROR WHILE SENDING NOTIFICATION TO GCM: %v\n", err)
			}
		}
	}()
}
