package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gitlab.com/fr5270937/notifications_service/internal/model"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/business_errors"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/validation"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/client"
)

type usecase struct {
	clientRepository desc.Repository
}

func NewClientUsecase(clientRepository desc.Repository) desc.Usecase {
	return usecase{clientRepository: clientRepository}
}

func (u usecase) CreateClient(ctx context.Context, req model.CreateClientRequest) (model.Client, error) {
	err := validation.ValidateCreateClientRequest(req)
	if err != nil {
		return model.Client{}, fmt.Errorf("validateCreateClientRequest: %w", err)
	}

	client, err := u.clientRepository.CreateClient(ctx, req)
	if err != nil {
		return model.Client{}, fmt.Errorf("clientRepository.CreateClient: %w", err)
	}

	return client, nil
}

func (u usecase) UpdateClient(ctx context.Context, id string, req model.UpdateClientRequest) (model.Client, error) {
	err := validation.ValidateUpdateClientRequest(id, req)
	if err != nil {
		return model.Client{}, fmt.Errorf("validateUpdateClientRequest: %w", err)
	}

	client, err := u.clientRepository.UpdateClient(ctx, id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Client{}, fmt.Errorf("clientRepository.UpdateClient: %w", business_errors.ErrClientNotFound)
		}

		return model.Client{}, fmt.Errorf("clientRepository.UpdateClient: %w", err)
	}

	return client, nil
}

func (u usecase) DeleteClient(ctx context.Context, id string) error {
	err := validation.ValidateDeleteClientRequest(id)
	if err != nil {
		return fmt.Errorf("validateDeleteClientRequest: %w", err)
	}

	err = u.clientRepository.DeleteClient(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("clientRepository.DeleteClient: %w", business_errors.ErrClientNotFound)
		}

		return fmt.Errorf("clientRepository.DeleteClient: %w", err)
	}

	return nil
}
