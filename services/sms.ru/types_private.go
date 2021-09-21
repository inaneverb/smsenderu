// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package smsenderu_smsru

import (
	"time"

	"github.com/qioalice/smsenderu"

	"github.com/qioalice/ekago/v3/ekatime"
)

func sendMessageRequestIsValid(req *smsenderu.SendMessageRequest) bool {

	isValid := req != nil &&
		(req.Recipient != "" ||
			(len(req.Recipients) > 0 && req.Recipients[0] != "")) &&
		req.Message != ""

	if !isValid {
		return false
	}

	if req.SendAt != 0 {
		if req.SendAt > ekatime.OnceInMinute.Now()+ekatime.SECONDS_IN_DAY*30 {
			return false
		}
		if req.SendAt <= ekatime.OnceInMinute.Now() {
			req.SendAt = 0
		}
	}

	if req.TTL < 1*time.Minute || req.TTL > 24*time.Hour {
		return false
	}

	return true
}

func sendMessageRequestWhyInvalid(req *smsenderu.SendMessageRequest) string {
	switch {
	case req == nil:
		return "Request object is nil"
	case req.Recipient == "" && (len(req.Recipients) == 0 || req.Recipients[0] == ""):
		return "No recipient is specified"
	case req.Message == "":
		return "No message body is specified"
	case req.SendAt != 0 && req.SendAt > ekatime.OnceInMinute.Now()+
		ekatime.SECONDS_IN_DAY*30:
		return "SendAt is more than 2 month over today"
	case req.TTL < 1*time.Minute || req.TTL > 24*time.Hour:
		return "TTL is in incorrect range, only [1m..24h] is allowed"
	default:
		return "Internal error. sendMessageRequestWhyInvalid()."
	}
}
