package mailing

import (
	"context"

	"gitlab.com/fr5270937/notifications_service/internal/model"
)

type Usecase interface {
	CreateMailing(ctx context.Context, req model.CreateMailingRequest) (model.Mailing, error)
	UpdateMailing(ctx context.Context, id string, req model.UpdateMailingRequest) (model.Mailing, error)
	DeleteMailing(ctx context.Context, id string) error
}
