// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package smsenderu_smsru_test

import (
	"testing"

	"github.com/qioalice/ekago/v2/ekalog"

	"github.com/qioalice/smsenderu"
	"github.com/qioalice/smsenderu/services/sms.ru"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
//==============================================================================//
	TOKEN = `171E5FA2-8D02-D511-1E8C-A6D6A6B984EB`
//==============================================================================//
)

func TestSenderSmsRu_Check(t *testing.T) {
	q := smsenderu_smsru.NewSender(TOKEN)
	err := q.Check()
	err.LogAsError("Failed to check SMS.RU credentials.")
	require.True(t, err.IsNil())
}

func TestSenderSmsRu_Balance(t *testing.T) {
//==============================================================================//
	expectedBalance := decimal.New(10, 0)
	//expectedBalance := decimal.NewFromString("10.0")
//==============================================================================//
	q := smsenderu_smsru.NewSender(TOKEN)
	balance, currency, err := q.Balance()
	err.LogAsError("Failed to check SMS.RU balance.")
	require.True(t, err.IsNil())
	require.EqualValues(t, "RUB", currency)
	require.EqualValues(t, expectedBalance.String(), balance.String())
}

func TestSenderSmsRu_Senders(t *testing.T) {
//==============================================================================//
	expectedSenders := []string{"< place your phone number here >"} // <--- EDIT ME
//==============================================================================//
	q := smsenderu_smsru.NewSender(TOKEN)
	senders, err := q.Senders()
	err.LogAsError("Failed to get SMS.RU senders.")
	require.True(t, err.IsNil())
	require.EqualValues(t, expectedSenders, senders)
}

func TestSenderSmsRu_Send(t *testing.T) {
//==============================================================================//
	req := &smsenderu.SendMessageRequest{
		Recipient: "< place your phone number here >", // <--- EDIT ME
		Message:   "UTF-8 日本語 тест\nCode: 1234",
	}
//==============================================================================//
	q := smsenderu_smsru.NewSender(TOKEN)
	resp, err := q.Send(req)
	err.LogAsError("Failed to send a message using SMS.RU.")
	require.True(t, err.IsNil())
	require.NotNil(t, resp)
	require.Len(t, resp.IDs, 1)
	require.Len(t, resp.IDs, 1)
	require.EqualValues(t, smsenderu_smsru.STATUS_OK, resp.ErrorCodes[0])
	ekalog.Debug("Send SMS ID: %s", resp.IDs[0])
}

func TestSenderSmsRu_Cost(t *testing.T) {
//==============================================================================//
	req := &smsenderu.SendMessageRequest{
		Recipient: "< place your phone number here >", // <--- EDIT ME
		Message:   "UTF-8 日本語 тест\nCode: 1234",
	}
//==============================================================================//
	q := smsenderu_smsru.NewSender(TOKEN)
	resp, err := q.Cost(req)
	err.LogAsError("Failed to get a cost of sending message(s) using SMS.RU.")
	require.True(t, err.IsNil())
	require.NotNil(t, resp)
	ekalog.Debug("Cost of SMS sending %s RUB", resp.Total)
}

func TestSenderSmsRu_Status(t *testing.T) {
//==============================================================================//
	sentSmsId := "202041-1000004"
//==============================================================================//
	q := smsenderu_smsru.NewSender(TOKEN)
	resp, err := q.Status(sentSmsId)
	err.LogAsError("Failed to get an info about sent message using SMS.RU.")
	require.True(t, err.IsNil())
	require.NotNil(t, resp)
	ekalog.Debug("Message info: %s", spew.Sdump(resp))
}
