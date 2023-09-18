package model

import (
	"time"
)

type Client struct {
	ID                 string    `json:"id"                   db:"id"                  ` // уникальный id клиента
	PhoneNumber        string    `json:"phone_number"         db:"phone_number"        ` // номер телефона клиента в формате 7XXXXXXXXXX (X - цифра от 0 до 9)
	MobileOperatorCode int32     `json:"mobile_operator_code" db:"mobile_operator_code"` // код мобильного оператора
	Tag                string    `json:"tag"                  db:"tag"                 ` // тег (произвольная метка)
	TimeZone           time.Time `json:"time_zone"            db:"time_zone"           ` // часовой пояс
}

type CreateClientRequest struct {
	PhoneNumber        string    `json:"phone_number"        `
	MobileOperatorCode int32     `json:"mobile_operator_code"`
	Tag                string    `json:"tag"                 `
	TimeZone           time.Time `json:"time_zone"           `
}

type UpdateClientRequest struct {
	PhoneNumber        string    `json:"phone_number"         `
	MobileOperatorCode int32     `json:"mobile_operator_code" `
	Tag                string    `json:"tag"                  `
	TimeZone           time.Time `json:"time_zone"            `
}
