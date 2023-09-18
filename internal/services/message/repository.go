package message

import (
	"context"

	"gitlab.com/fr5270937/notifications_service/internal/model"
)

type Repository interface {
	CreateMessage(ctx context.Context, message model.Message) (model.Message, error)
	UpdateMessage(ctx context.Context, id int64, message model.Message) (model.Message, error)
	GetMessagesByMailingID(ctx context.Context, id string) ([]model.Message, error)
	GetStatusMessagesByMailingIDAndStatus(ctx context.Context, id string, status int32) ([]model.Message, error)
}
