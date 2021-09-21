// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package smsenderu_smsru_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/qioalice/ekago/v3/ekalog"

	"github.com/qioalice/smsenderu"
	"github.com/qioalice/smsenderu/services/sms.ru"
)

const (
	//==============================================================================//
	TOKEN = `< place your token here >`

//==============================================================================//
)

func TestSenderSmsRu_Check(t *testing.T) {
	q := smsenderu_smsru.NewSender(TOKEN)
	err := q.Check()
	ekalog.Errore("Failed to check SMS.RU credentials.", err)
	require.True(t, err.IsNil())
}

func TestSenderSmsRu_Balance(t *testing.T) {
	//==============================================================================//
	expectedBalance := decimal.New(10, 0)
	//expectedBalance := decimal.NewFromString("10.0")
	//==============================================================================//
	q := smsenderu_smsru.NewSender(TOKEN)
	balance, currency, err := q.Balance()
	ekalog.Errore("Failed to check SMS.RU balance.", err)
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
	ekalog.Errore("Failed to get SMS.RU senders.", err)
	require.True(t, err.IsNil())
	spew.Dump(senders)
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
	ekalog.Errore("Failed to send a message using SMS.RU.", err)
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
	ekalog.Errore("Failed to get a cost of sending message(s) using SMS.RU.", err)
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
	ekalog.Errore("Failed to get an info about sent message using SMS.RU.", err)
	require.True(t, err.IsNil())
	require.NotNil(t, resp)
	ekalog.Debug("Message info: %s", spew.Sdump(resp))
}
