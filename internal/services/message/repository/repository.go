package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/message"
)

type repository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) desc.Repository {
	return repository{db: db}
}

func (r repository) CreateMessage(ctx context.Context, message model.Message) (model.Message, error) {
	err := r.db.QueryRowContext(ctx, `INSERT INTO message (status, mailing_id, receiver_client_id) VALUES ($1, $2, $3) RETURNING *`,
		message.Status, message.MailingID, message.ReceiverClientID).
		Scan(&message.ID, &message.CreatedAt, &message.Status, &message.MailingID, &message.ReceiverClientID)
	if err != nil {
		return model.Message{}, err
	}

	return message, nil
}

func (r repository) UpdateMessage(ctx context.Context, id int64, message model.Message) (model.Message, error) {
	var t time.Time

	builder := squirrel.Update("message").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *")
	if message.CreatedAt != t {
		builder = builder.Set("created_at", message.CreatedAt)
	}
	if message.Status != 0 {
		builder = builder.Set("status", message.Status)
	}
	if message.MailingID != "" {
		builder = builder.Set("mailing_id", message.MailingID)
	}
	if message.ReceiverClientID != "" {
		builder = builder.Set("receiver_client_id", message.ReceiverClientID)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Message{}, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&message.ID, &message.CreatedAt, &message.Status, &message.MailingID, &message.ReceiverClientID)
	if err != nil {
		return model.Message{}, err
	}

	return message, nil
}

func (r repository) GetMessagesByMailingID(ctx context.Context, id string) ([]model.Message, error) {
	var messages []model.Message

	err := r.db.SelectContext(ctx, &messages, "SELECT * FROM message WHERE mailing_id=$1 GROUP BY message.id, status", id)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r repository) GetStatusMessagesByMailingIDAndStatus(ctx context.Context, id string, status int32) ([]model.Message, error) {
	var messages []model.Message

	err := r.db.SelectContext(ctx, &messages, "SELECT * FROM message WHERE mailing_id=$1 AND status=$2", id, status)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
