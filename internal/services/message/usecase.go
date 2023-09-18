package message

import (
	"context"

	"gitlab.com/fr5270937/notifications_service/internal/model"
)

type Usecase interface {
	StartMailing(ctx context.Context, mailing model.Mailing) error
	GetStatisticMessages(ctx context.Context) ([]model.CommonStatisticMessages, error)
	GetDetailStatisticMessagesByMailing(ctx context.Context, id string) ([]model.Message, error)
	HandlingActiveMessages(ctx context.Context) error
}
