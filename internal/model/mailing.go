package model

import (
	"time"
)

type Mailing struct {
	ID         string    `json:"id"          db:"id"         ` // уникальный id рассылки
	CreatedAt  time.Time `json:"created_at"  db:"created_at" ` // дата и время запуска рассылки
	Text       string    `json:"text"        db:"text"       ` // текст сообщения для доставки клиенту
	Filter     string    `json:"filter"      db:"filter"     ` // фильтр свойств клиентов, на которых должна быть произведена рассылка (код мобильного оператора, тег)
	FinishedAt time.Time `json:"finished_at" db:"finished_at"` // дата и время окончания рассылки: если по каким-то причинам не успели разослать все сообщения - никакие сообщения клиентам после этого времени доставляться не должны
}

type PendingMailing struct {
	ID        string // уникальный id отложенной рассылки
	MailingID string // id рассылки из общей таблицы
	Status    int    // статус отложенной рассылки
}

type CreateMailingRequest struct {
	CreatedAt  time.Time `json:"created_at" `
	Text       string    `json:"text"       `
	Filter     string    `json:"filter"     `
	FinishedAt time.Time `json:"finished_at"`
}

type UpdateMailingRequest struct {
	CreatedAt  time.Time `json:"created_at" `
	Text       string    `json:"text"       `
	Filter     string    `json:"filter"     `
	FinishedAt time.Time `json:"finished_at"`
}

type UpdatePendingMailing struct {
	MailingID string
	Status    int
}

// PendingMailingStatus
const (
	PENDING_MAILING_STATUS_UNKNOWN = 0
	PENDING_MAILING_STATUS_CREATED = 1
	PENDING_MAILING_STATUS_IN_WORK = 2
)
