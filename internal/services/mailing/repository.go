package mailing

import (
	"context"

	"gitlab.com/fr5270937/notifications_service/internal/model"
)

type Repository interface {
	CreateMailing(ctx context.Context, req model.CreateMailingRequest) (model.Mailing, error)
	UpdateMailing(ctx context.Context, id string, req model.UpdateMailingRequest) (model.Mailing, error)
	DeleteMailing(ctx context.Context, id string) error
	UpdatePendingMailing(ctx context.Context, mailingID string, mailing model.UpdatePendingMailing) (model.PendingMailing, error)
	DeletePendingMailing(ctx context.Context, id string) error
	GetAllMailings(ctx context.Context) ([]model.Mailing, error)
	GetActiveMailings(ctx context.Context) ([]model.Mailing, error)
	GetMailingsToSending(ctx context.Context) ([]model.Mailing, error)
}
