// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package smsenderu

import (
	"github.com/shopspring/decimal"

	"github.com/qioalice/ekago/v3/ekaerr"
)

type (
	Sender interface {

		// Check must checks whether Sender was initialized properly
		// and it's very recommend to do a test API service request (kinda ping)
		// or another way attempt to figure out whether provided service's credentials
		// (login-password, login-token, something else) are valid.
		Check() *ekaerr.Error

		Balance() (balance decimal.Decimal, currency string, err *ekaerr.Error)
		BalanceIn(currency string) (balance decimal.Decimal, err *ekaerr.Error)

		// Senders must return a list of available senders,
		// that can be used through API as a displayable sender for recipients.
		Senders() ([]string, *ekaerr.Error)

		Send(req *SendMessageRequest) (resp *SendMessageResponse, err *ekaerr.Error)
		Cost(req *SendMessageRequest) (resp *CostSendMessageResponse, err *ekaerr.Error)

		Status(sentSmsId string) (resp *StatusMessageResponse, err *ekaerr.Error)
	}
)
