package client

import (
	"context"

	"gitlab.com/fr5270937/notifications_service/internal/model"
)

type Usecase interface {
	CreateClient(ctx context.Context, req model.CreateClientRequest) (model.Client, error)
	UpdateClient(ctx context.Context, id string, req model.UpdateClientRequest) (model.Client, error)
	DeleteClient(ctx context.Context, id string) error
}
