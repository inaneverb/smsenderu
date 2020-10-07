// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package smsenderu_smsru

import (
	"strconv"
	"strings"

	"github.com/qioalice/ekago/v2/ekaerr"

	"github.com/qioalice/smsenderu"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

func NewSender(token string) smsenderu.Sender {
	return &senderSmsRu{ token: token }
}

func (q *senderSmsRu) Check() *ekaerr.Error {
	// https://sms.ru/api/auth_check
	const s = "SMS.RU: Failed to check whether provided API token is valid. "
	switch {

	case q == nil:
		return ekaerr.IllegalArgument.
			New(s + "Invalid sender object. Did you use NewSender() constructor correctly?").
			Throw()

	case q.token == "":
		return ekaerr.IllegalArgument.
			New(s + "API token is not provided or empty.").
			Throw()
	}

	const URL = "https://sms.ru/auth/check"
	fhReq := fasthttp.AcquireRequest(); fasthttp.ReleaseRequest(fhReq)
	fhResp := fasthttp.AcquireResponse(); fasthttp.ReleaseResponse(fhResp)

	fhReq.SetRequestURI(URL)
	fhReq.URI().QueryArgs().Add("api_id", q.token)

	_, err := q.do(fhReq, fhResp, 0)
	if err.IsNotNil() {
		return err.
			AddMessage(s).
			Throw()
	}

	return nil
}

func (q *senderSmsRu) Balance() (decimal.Decimal, string, *ekaerr.Error) {
	// https://sms.ru/api/balance
	const s = "SMS.RU: Failed to get balance. "
	switch {

	case q == nil:
		return decimal.Zero, "", ekaerr.IllegalArgument.
			New(s + "Invalid sender object. Did you use NewSender() constructor correctly?").
			Throw()

	case q.token == "":
		return decimal.Zero, "", ekaerr.IllegalArgument.
			New(s + "API token is not provided or empty.").
			Throw()
	}

	const URL = "https://sms.ru/my/balance"
	fhReq := fasthttp.AcquireRequest(); fasthttp.ReleaseRequest(fhReq)
	fhResp := fasthttp.AcquireResponse(); fasthttp.ReleaseResponse(fhResp)

	fhReq.SetRequestURI(URL)
	fhReq.URI().QueryArgs().Add("api_id", q.token)

	respParts, err := q.do(fhReq, fhResp, 1)
	if err.IsNotNil() {
		return decimal.Zero, "", err.
			AddMessage(s).
			Throw()
	}

	balance, legacyErr := decimal.NewFromString(string(respParts[0]))
	if legacyErr != nil {
		return decimal.Zero, "", ekaerr.IllegalFormat.
			Wrap(legacyErr, s + "Cannot decode balance value from API response.").
			AddFields("smsru_balance_raw", string(respParts[0])).
			Throw()
	}

	return balance, "RUB", nil
}

func (q *senderSmsRu) BalanceIn(currency string) (balance decimal.Decimal, err *ekaerr.Error) {
	const s = "SMS.RU: Failed to get balance in the specified currency. "

	currency = strings.TrimSpace(currency)
	currency = strings.ToUpper(currency)

	switch currency {

	case "RUB":
		balance, _, err = q.Balance()
		return balance, err.
			AddMessage(s).
			Throw()

	case "":
		return decimal.Zero, ekaerr.IllegalArgument.
			New(s + "Currency is empty or not provided.").
			Throw()

	default:
		return decimal.Zero, ekaerr.IllegalArgument.
			New(s + "Incorrect currency. Only 'RUB' is supported.").
			AddFields("smsru_required_currency", currency).
			Throw()
	}
}

func (q *senderSmsRu) Senders() ([]string, *ekaerr.Error) {
	// https://sms.ru/api/senders
	const s = "SMS.RU: Failed to get registered senders. "
	switch {

	case q == nil:
		return nil, ekaerr.IllegalArgument.
			New(s + "Invalid sender object. Did you use NewSender() constructor correctly?").
			Throw()

	case q.token == "":
		return nil, ekaerr.IllegalArgument.
			New(s + "API token is not provided or empty.").
			Throw()
	}

	const URL = "https://sms.ru/my/senders"
	fhReq := fasthttp.AcquireRequest(); fasthttp.ReleaseRequest(fhReq)
	fhResp := fasthttp.AcquireResponse(); fasthttp.ReleaseResponse(fhResp)

	fhReq.SetRequestURI(URL)
	fhReq.URI().QueryArgs().Add("api_id", q.token)

	respParts, err := q.do(fhReq, fhResp, 0)
	if err.IsNotNil() {
		return nil, err.
			AddMessage(s).
			Throw()
	}

	if len(respParts) == 0 {
		return nil, nil
	}

	senders := make([]string, len(respParts))
	for i, n := 0, len(respParts); i < n; i++ {
		senders[i] = string(respParts[i])
	}

	return senders, nil
}

func (q *senderSmsRu) Send(

	req *smsenderu.SendMessageRequest,
) (
	resp *smsenderu.SendMessageResponse,
	err  *ekaerr.Error,
) {
	// https://sms.ru/api/send
	const s = "SMS.RU: Failed to send a message(s). "
	switch {

	case q == nil:
		return nil, ekaerr.IllegalArgument.
			New(s + "Invalid sender object. Did you use NewSender() constructor correctly?").
			Throw()

	case q.token == "":
		return nil, ekaerr.IllegalArgument.
			New(s + "API token is not provided or empty.").
			Throw()

	case !sendMessageRequestIsValid(req):
		return nil, ekaerr.IllegalArgument.
			New(s + "Incorrect argument(s) of sending message request.").
			AddFields(
				"smsru_send_request_why_invalid", sendMessageRequestWhyInvalid(req),
				"smsru_send_request_dump", spew.Sdump(req)).
			Throw()
	}

	const URL = "https://sms.ru/sms/send"
	fhReq := fasthttp.AcquireRequest(); fasthttp.ReleaseRequest(fhReq)
	fhResp := fasthttp.AcquireResponse(); fasthttp.ReleaseResponse(fhResp)

	fhReq.SetRequestURI(URL)
	fhReq.URI().QueryArgs().Add("api_id", q.token)

	if req.Recipient != "" {
		req.Recipients = []string{req.Recipient}
	}

	fhReq.URI().QueryArgs().Add("to", strings.Join(req.Recipients, ","))
	fhReq.URI().QueryArgs().Add("msg", req.Message)

	if req.From != "" {
		fhReq.URI().QueryArgs().Add("from", req.From)
	}

	if req.UserIP != "" {
		fhReq.URI().QueryArgs().Add("ip", req.UserIP)
	}

	if req.SendAt != 0 {
		v := strconv.FormatInt(req.SendAt.I64(), 10)
		fhReq.URI().QueryArgs().Add("time", v)
	}

	if req.TTL != 0 {
		v := strconv.Itoa(int(req.TTL.Minutes()))
		fhReq.URI().QueryArgs().Add("ttl", v)
	}

	if req.EnableUserLocation {
		fhReq.URI().QueryArgs().Add("daytime", "1")
	}

	if req.DoTransliterate {
		fhReq.URI().QueryArgs().Add("translit", "1")
	}

	// Response contain an info about each phone number's sending message
	// +1 more row (the current balance after sending).

	var respParts [][]byte
	respParts, err = q.do(fhReq, fhResp, len(req.Recipients)+1)
	if err.IsNotNil() {
		return nil, err.
			AddMessage(s).
			Throw()
	}

	if len(respParts) != len(req.Recipients) +1 {
		return nil, ekaerr.IllegalFormat.
			New(s + "Mismatch sizes of the response and request phone numbers.").
			AddFields(
				"smsru_send_request_phone_numbers", len(req.Recipients),
				"smsru_send_response_phone_numbers", len(respParts),
				"smsru_send_request_dump", spew.Sdump(req),
				"smsru_send_response_raw", string(fhResp.Body())).
			Throw()
	}

	resp = &smsenderu.SendMessageResponse{
		IDs:        make([]string, len(req.Recipients)),
		ErrorCodes: make([]int, len(req.Recipients)),
	}

	// It's unnecessary but would even caller use SendMessageRequest object
	// after that call?
	if req.Recipient != "" {
		req.Recipients = nil
	}

	for i, n := 0, len(respParts)-1; i < n; i++ {
		if len(respParts[i]) <= 3 {
			// Seems like error code
			resp.ErrorCodes[i], _ = strconv.Atoi(string(respParts[i]))
		} else {
			resp.IDs[i]        = string(respParts[i])
			resp.ErrorCodes[i] = STATUS_OK
		}
	}

	return resp, nil
}

func (q *senderSmsRu) Cost(

	req *smsenderu.SendMessageRequest,
) (
	resp *smsenderu.CostSendMessageResponse,
	err  *ekaerr.Error,
) {
	// https://sms.ru/api/cost
	const s = "SMS.RU: Failed to get an info about cost of sending a message(s). "
	switch {

	case q == nil:
		return nil, ekaerr.IllegalArgument.
			New(s + "Invalid sender object. Did you use NewSender() constructor correctly?").
			Throw()

	case q.token == "":
		return nil, ekaerr.IllegalArgument.
			New(s + "API token is not provided or empty.").
			Throw()

	case !sendMessageRequestIsValid(req):
		return nil, ekaerr.IllegalArgument.
			New(s + "Incorrect argument(s) of sending message request.").
			AddFields(
				"smsru_cost_request_why_invalid", sendMessageRequestWhyInvalid(req),
				"smsru_cost_request_dump", spew.Sdump(req)).
			Throw()
	}

	const URL = "https://sms.ru/sms/cost"
	fhReq := fasthttp.AcquireRequest(); fasthttp.ReleaseRequest(fhReq)
	fhResp := fasthttp.AcquireResponse(); fasthttp.ReleaseResponse(fhResp)

	fhReq.SetRequestURI(URL)
	fhReq.URI().QueryArgs().Add("api_id", q.token)

	if req.Recipient != "" {
		req.Recipients = []string{req.Recipient}
	}

	fhReq.URI().QueryArgs().Add("to", strings.Join(req.Recipients, ","))
	fhReq.URI().QueryArgs().Add("msg", req.Message)

	if req.From != "" {
		fhReq.URI().QueryArgs().Add("from", req.From)
	}

	if req.DoTransliterate {
		fhReq.URI().QueryArgs().Add("translit", "1")
	}

	// Response contain an info about each phone number's sending message
	// +1 more row (the current balance after sending).

	var respParts [][]byte
	respParts, err = q.do(fhReq, fhResp, 2)
	if err.IsNotNil() {
		return nil, err.
			AddMessage(s).
			Throw()
	}

	resp = &smsenderu.CostSendMessageResponse{
		Costs: nil, // not supported by https://sms.ru/
		Total: decimal.Zero,
	}

	// It's unnecessary but would even caller use SendMessageRequest object
	// after that call?
	if req.Recipient != "" {
		req.Recipients = nil
	}

	cost, legacyErr := decimal.NewFromString(string(respParts[0]))
	if legacyErr != nil {
		return nil, ekaerr.IllegalFormat.
			Wrap(legacyErr, s + "Cannot decode the cost of sending message(s).").
			AddFields(
				"smsru_cost_request_dump", spew.Sdump(req),
				"smsru_cost_response_raw", string(fhResp.Body())).
			Throw()
	}

	resp.Total = cost
	return resp, nil
}

func (q *senderSmsRu) Status(

	sentSmsId string,
) (
	resp *smsenderu.StatusMessageResponse,
	err *ekaerr.Error,
) {
	// https://sms.ru/api/status
	const s = "SMS.RU: Failed to get an info about message. "
	switch {

	case q == nil:
		return nil, ekaerr.IllegalArgument.
			New(s + "Invalid sender object. Did you use NewSender() constructor correctly?").
			Throw()

	case q.token == "":
		return nil, ekaerr.IllegalArgument.
			New(s + "API token is not provided or empty.").
			Throw()

	case sentSmsId == "":
		return nil, ekaerr.IllegalArgument.
			New(s + "SMS ID is empty or not provided.").
			Throw()
	}

	const URL = "https://sms.ru/sms/status"
	fhReq := fasthttp.AcquireRequest(); fasthttp.ReleaseRequest(fhReq)
	fhResp := fasthttp.AcquireResponse(); fasthttp.ReleaseResponse(fhResp)

	fhReq.SetRequestURI(URL)
	fhReq.URI().QueryArgs().Add("api_id", q.token)

	fhReq.URI().QueryArgs().Add("sms_id", sentSmsId)

	var respParts [][]byte
	respParts, err = q.do(fhReq, fhResp, 1)
	if err.IsNotNil() {
		return nil, err.
			AddMessage(s).
			Throw()
	}

	resp = &smsenderu.StatusMessageResponse{
		ID: sentSmsId,
	}

	resp.ErrorCode, _ = strconv.Atoi(string(respParts[0]))

	return resp, nil
}
