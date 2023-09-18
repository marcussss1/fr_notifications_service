package client

import (
	"context"

	"gitlab.com/fr5270937/notifications_service/internal/model"
)

type Repository interface {
	CreateClient(ctx context.Context, req model.CreateClientRequest) (model.Client, error)
	UpdateClient(ctx context.Context, id string, req model.UpdateClientRequest) (model.Client, error)
	DeleteClient(ctx context.Context, id string) error
	GetClientsFromMobileOperatorCode(ctx context.Context, mobileOperatorCode int32) ([]model.Client, error)
	GetClientsFromTag(ctx context.Context, tag string) ([]model.Client, error)
	GetClientFromID(ctx context.Context, id string) (model.Client, error)
}
