// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package smsenderu_smsru

import (
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/qioalice/ekago/v3/ekaerr"
	"github.com/qioalice/ekago/v3/ekastr"
)

type (
	senderSmsRu struct {
		fhc   fasthttp.Client
		token string
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	ERROR_CODE_MESSAGE_NOT_FOUND = -1

	STATUS_OK                                 = 100
	STATUS_PENDING_BY_OPERATOR                = 101
	STATUS_PENDING                            = 102
	STATUS_DELIVERED                          = 103
	STATUS_NOT_DELIVERED_TIMEOUT              = 104
	STATUS_NOT_DELIVERED_REJECTED_BY_OPERATOR = 105
	STATUS_NOT_DELIVERED_PHONE_FAILURE        = 106
	STATUS_NOT_DELIVERED_UNKNOWN              = 107
	STATUS_NOT_DELIVERED_REJECTED             = 108
	STATUS_NOT_DELIVERED_BAD_ROUTE            = 150
	STATUS_READ                               = 110

	ERROR_CODE_INCORRECT_API_TOKEN                                        = 200
	ERROR_CODE_NOT_ENOUGH_MONEY                                           = 201
	ERROR_CODE_BAD_PHONE_NUMBER                                           = 202
	ERROR_CODE_NO_MESSAGE_BODY                                            = 203
	ERROR_CODE_SENDER_IS_NOT_APPROVED                                     = 204
	ERROR_CODE_MESSAGE_BODY_TOO_LARGE                                     = 205
	ERROR_CODE_USER_DEFINED_LIMIT_IS_REACHED                              = 206
	ERROR_CODE_BAD_ROUTE                                                  = 207
	ERROR_CODE_INCORRECT_TIME                                             = 208
	ERROR_CODE_PHONE_NUMBER_IS_BLOCKED_BY_USER                            = 209
	ERROR_CODE_HTTP_METHOD_IS_NOT_ALLOWED                                 = 210
	ERROR_CODE_API_METHOD_NOT_FOUND                                       = 211
	ERROR_CODE_INCORRECT_MESSAGE_BODY_ENCODING                            = 212
	ERROR_CODE_TOO_MUCH_PHONE_NUMBERS                                     = 213
	ERROR_CODE_TEMPORARY_UNAVAILABLE                                      = 220
	ERROR_CODE_DAILY_PER_PHONE_NUMBER_LIMIT_IS_REACHED                    = 230
	ERROR_CODE_SAME_MESSAGES_PER_MINUTE_PER_PHONE_NUMBER_LIMIT_IS_REACHED = 231
	ERROR_CODE_SAME_MESSAGES_PER_DAY_PER_PHONE_NUMBER_LIMIT_IS_REACHED    = 232
	ERROR_CODE_SPAM_DETECTED                                              = 233

	ERROR_CODE_EXPIRED_API_TOKEN                     = 300
	ERROR_CODE_INCORRECT_LOGIN_OR_PASSWORD           = 301
	ERROR_CODE_AUTHORIZED_BUT_NOT_ACTIVATED          = 302
	ERROR_CODE_AUTHORIZED_BUT_2FA_INCORRECT          = 303
	ERROR_CODE_AUTHORIZED_BUT_TOO_MUCH_2FA_SENT      = 304
	ERROR_CODE_AUTHORIZED_BUT_TOO_MUCH_INCORRECT_2FA = 305

	ERROR_CODE_INTERNAL_SERVER_ERROR = 500

	ERROR_CODE_CALLBACK_INCORRECT_URL = 901
	ERROR_CODE_CALLBACK_NOT_FOUND     = 902
)

var (
	statusCodeMeaningMap = map[int]string{
		ERROR_CODE_MESSAGE_NOT_FOUND:                       "Message not found",
		STATUS_OK:                                          "OK",
		STATUS_PENDING_BY_OPERATOR:                         "Sent, waiting for operator",
		STATUS_PENDING:                                     "Sent, waiting for delivery",
		STATUS_DELIVERED:                                   "Delivered",
		STATUS_NOT_DELIVERED_TIMEOUT:                       "Not delivered: Timeout",
		STATUS_NOT_DELIVERED_REJECTED_BY_OPERATOR:          "Not delivered: Rejected by operator",
		STATUS_NOT_DELIVERED_PHONE_FAILURE:                 "Not delivered: Phone failure",
		STATUS_NOT_DELIVERED_UNKNOWN:                       "Not delivered: Unknown error",
		STATUS_NOT_DELIVERED_REJECTED:                      "Not delivered: Rejected",
		STATUS_NOT_DELIVERED_BAD_ROUTE:                     "Not delivered: Bad route",
		STATUS_READ:                                        "Message has been read",
		ERROR_CODE_INCORRECT_API_TOKEN:                     "Bad request: Incorrect API token",
		ERROR_CODE_NOT_ENOUGH_MONEY:                        "Bad request: Not enough money",
		ERROR_CODE_BAD_PHONE_NUMBER:                        "Bad request: Incorrect phone number",
		ERROR_CODE_NO_MESSAGE_BODY:                         "Bad request: No message body",
		ERROR_CODE_SENDER_IS_NOT_APPROVED:                  "Bad request: Sender is not approved",
		ERROR_CODE_MESSAGE_BODY_TOO_LARGE:                  "Bad request: Body too large",
		ERROR_CODE_USER_DEFINED_LIMIT_IS_REACHED:           "Limits: Admin defined limit is reached",
		ERROR_CODE_BAD_ROUTE:                               "Bad request: Bad route",
		ERROR_CODE_INCORRECT_TIME:                          "Bad request: Incorrect time",
		ERROR_CODE_PHONE_NUMBER_IS_BLOCKED_BY_USER:         "Bad request: Phone number is locked by admin",
		ERROR_CODE_HTTP_METHOD_IS_NOT_ALLOWED:              "Bad request: HTTP method not allowed",
		ERROR_CODE_API_METHOD_NOT_FOUND:                    "Bad request: HTTP route not found",
		ERROR_CODE_INCORRECT_MESSAGE_BODY_ENCODING:         "Bad request: Incorrect body encoding",
		ERROR_CODE_TOO_MUCH_PHONE_NUMBERS:                  "Bad request: Too much phone numbers (recipients)",
		ERROR_CODE_TEMPORARY_UNAVAILABLE:                   "Server: Temporary unavailable",
		ERROR_CODE_DAILY_PER_PHONE_NUMBER_LIMIT_IS_REACHED: "Limits: Daily limit per phone number is reached",
		ERROR_CODE_SAME_MESSAGES_PER_MINUTE_PER_PHONE_NUMBER_LIMIT_IS_REACHED: "Limits: Same message per minute per phone number is reached",
		ERROR_CODE_SAME_MESSAGES_PER_DAY_PER_PHONE_NUMBER_LIMIT_IS_REACHED:    "Limits: Same message per day per phone number is reached",
		ERROR_CODE_SPAM_DETECTED:                         "Limits: Spam detected",
		ERROR_CODE_EXPIRED_API_TOKEN:                     "Bad request: API token is expired",
		ERROR_CODE_INCORRECT_LOGIN_OR_PASSWORD:           "Bad request: Incorrect login or password",
		ERROR_CODE_AUTHORIZED_BUT_NOT_ACTIVATED:          "Bad request: Authorized, but not activated",
		ERROR_CODE_AUTHORIZED_BUT_2FA_INCORRECT:          "Bad request: Authorized, but 2FA is incorrect",
		ERROR_CODE_AUTHORIZED_BUT_TOO_MUCH_2FA_SENT:      "Bad request: Authorized, but too much sending 2FA",
		ERROR_CODE_AUTHORIZED_BUT_TOO_MUCH_INCORRECT_2FA: "Bad request: Authorized, but too much incorrect 2FA",
		ERROR_CODE_INTERNAL_SERVER_ERROR:                 "Server: Internal server error",
		ERROR_CODE_CALLBACK_INCORRECT_URL:                "Callbacks: Incorrect URL",
		ERROR_CODE_CALLBACK_NOT_FOUND:                    "CallbacksL No registered callback",
	}
)

// Original:
// https://sms.ru/api/status
//
//-1	Сообщение не найдено
//
//100	Запрос выполнен или сообщение находится в нашей очереди
//101	Сообщение передается оператору
//102	Сообщение отправлено (в пути)
//103	Сообщение доставлено
//104	Не может быть доставлено: время жизни истекло
//105	Не может быть доставлено: удалено оператором
//106	Не может быть доставлено: сбой в телефоне
//107	Не может быть доставлено: неизвестная причина
//108	Не может быть доставлено: отклонено
//110	Сообщение прочитано (для Viber, временно не работает)
//150	Не может быть доставлено: не найден маршрут на данный номер
//
//200	Неправильный api_id
//201	Не хватает средств на лицевом счету
//202	Неправильно указан номер телефона получателя, либо на него нет маршрута
//203	Нет текста сообщения
//204	Имя отправителя не согласовано с администрацией
//205	Сообщение слишком длинное (превышает 8 СМС)
//206	Будет превышен или уже превышен дневной лимит на отправку сообщений
//207	На этот номер нет маршрута для доставки сообщений
//208	Параметр time указан неправильно
//209	Вы добавили этот номер (или один из номеров) в стоп-лист
//210	Используется GET, где необходимо использовать POST
//211	Метод не найден
//212	Текст сообщения необходимо передать в кодировке UTF-8 (вы передали в другой кодировке)
//213	Указано более 100 номеров в списке получателей
//220	Сервис временно недоступен, попробуйте чуть позже
//230	Превышен общий лимит количества сообщений на этот номер в день
//231	Превышен лимит одинаковых сообщений на этот номер в минуту
//232	Превышен лимит одинаковых сообщений на этот номер в день
//233	Превышен лимит отправки повторных сообщений с кодом на этот номер за короткий промежуток времени ("защита от мошенников", можно отключить в разделе "Настройки")
//
//300	Неправильный token (возможно истек срок действия, либо ваш IP изменился)
//301	Неправильный api_id, либо логин/пароль
//302	Пользователь авторизован, но аккаунт не подтвержден (пользователь не ввел код, присланный в регистрационной смс)
//303	Код подтверждения неверен
//304	Отправлено слишком много кодов подтверждения. Пожалуйста, повторите запрос позднее
//305	Слишком много неверных вводов кода, повторите попытку позднее
//
//500	Ошибка на сервере. Повторите запрос.
//
//901	Callback: URL неверный (не начинается на http://)
//902	Callback: Обработчик не найден (возможно был удален ранее)

func (q *senderSmsRu) do(

	fhReq *fasthttp.Request,
	fhResp *fasthttp.Response,
	requiredParts int,
) (
	parts [][]byte,
	err *ekaerr.Error,
) {
	legacyErr := q.fhc.DoRedirects(fhReq, fhResp, 5)
	if legacyErr != nil {
		return nil, ekaerr.ServiceUnavailable.
			Wrap(legacyErr, "Failed to perform HTTP request.").
			Throw()
	}

	if httpCode := fhResp.StatusCode(); httpCode != fasthttp.StatusOK {
		return nil, ekaerr.RejectedOperation.
			New("API response finished with other than HTTP 200 status code.").
			WithInt("smsru_response_http_code", httpCode).
			Throw()
	}

	var (
		statusCode        int
		statusCodeMeaning = "<UnknownStatus>"
	)

	statusCode, parts = q.decodeResponse(fhResp.Body())
	if statusCode == 0 {
		return nil, ekaerr.IllegalFormat.
			New("Failed to decode API response. It is empty or without status.").
			WithString("smsru_response_raw", ekastr.B2S(fhResp.Body())).
			Throw()
	}

	if statusCode != STATUS_OK {
		if statusCodeMeaning_, ok := statusCodeMeaningMap[statusCode]; ok {
			statusCodeMeaning = statusCodeMeaning_
		}
		return nil, ekaerr.IllegalFormat.
			New("API response finished with not OK code.").
			WithInt("smsru_response_status_code", statusCode).
			WithString("smsru_response_status_code_meaning", statusCodeMeaning).
			WithString("smsru_response_raw", ekastr.B2S(fhResp.Body())).
			Throw()
	}

	if len(parts) < requiredParts {
		return nil, ekaerr.IllegalFormat.
			New("Failed to decode API response. Unexpected number of parts.").
			WithInt("smsru_response_required_parts", requiredParts).
			WithInt("smsru_response_got_parts", len(parts)).
			WithString("smsru_response_raw", ekastr.B2S(fhResp.Body())).
			Throw()
	}

	return parts, nil
}

func (q *senderSmsRu) decodeResponse(b []byte) (statusCode int, parts [][]byte) {

	n := len(b)
	if n == 0 {
		return 0, nil
	}

	parts = make([][]byte, 0, 16)

	for i := 0; i < n; i++ {
		j := i
		for ; i < n && b[i] != '\n'; i++ {
		}
		parts = append(parts, b[j:i])
	}

	statusCodePart := parts[0]
	parts = parts[1:]

	statusCode, _ = strconv.Atoi(string(statusCodePart))

	return statusCode, parts
}
