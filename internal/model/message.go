package model

import (
	"time"
)

type Message struct {
	ID               int64     `json:"id"                 db:"id"                ` // уникальный id сообщения
	CreatedAt        time.Time `json:"created_at"         db:"created_at"        ` // дата и время создания (отправки)
	Status           int32     `json:"status"             db:"status"            ` // статус отправки
	MailingID        string    `json:"mailing_id"         db:"mailing_id"        ` // id рассылки, в рамках которой было отправлено сообщение
	ReceiverClientID string    `json:"receiver_client_id" db:"receiver_client_id"` // id клиента, которому отправили
}

type SendMessageRequest struct {
	ID          int64  `json:"id"   `
	PhoneNumber int64  `json:"phone"`
	Text        string `json:"text" `
}

type CommonStatisticMessages struct {
	Mailing  Mailing   `json:"mailing" `
	Messages []Message `json:"messages"`
}

// MessageStatus
const (
	MESSAGE_STATUS_UNKNOWN   = 0
	MESSAGE_STATUS_CREATED   = 1
	MESSAGE_STATUS_SUCCEEDED = 2
	MESSAGE_STATUS_FAILED    = 3
)
