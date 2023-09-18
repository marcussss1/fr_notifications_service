package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/business_errors"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/validation"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/mailing"
	"gitlab.com/fr5270937/notifications_service/internal/services/message"
)

type usecase struct {
	messageUsecase    message.Usecase
	mailingRepository desc.Repository
}

func NewMailingUsecase(mailingRepository desc.Repository, messageUsecase message.Usecase) desc.Usecase {
	return usecase{messageUsecase: messageUsecase, mailingRepository: mailingRepository}
}

func (u usecase) CreateMailing(ctx context.Context, req model.CreateMailingRequest) (model.Mailing, error) {
	var t time.Time

	err := validation.ValidateCreateMailingRequest(req)
	if err != nil {
		return model.Mailing{}, fmt.Errorf("validateCreateMailingRequest: %w", err)
	}
	if req.CreatedAt == t {
		req.CreatedAt = time.Now()
	}
	if req.FinishedAt == t {
		req.FinishedAt = time.Now().Add(1 * time.Minute)
	}

	timeNow := time.Now()
	mailing, err := u.mailingRepository.CreateMailing(ctx, req)
	if err != nil {
		return model.Mailing{}, fmt.Errorf("mailingRepository.CreateMailing: %w", err)
	}

	if timeNow.String() > mailing.CreatedAt.String() && timeNow.String() < mailing.FinishedAt.String() {
		go func() {
			errIn := u.messageUsecase.StartMailing(ctx, mailing)
			if errIn != nil {
				log.Error(fmt.Errorf("messageUsecase.StartMailing: %w", errIn))
			}
		}()
	}

	return mailing, nil
}

func (u usecase) UpdateMailing(ctx context.Context, id string, req model.UpdateMailingRequest) (model.Mailing, error) {
	err := validation.ValidateUpdateMailingRequest(id, req)
	if err != nil {
		return model.Mailing{}, fmt.Errorf("validateUpdateMailingRequest: %w", err)
	}

	client, err := u.mailingRepository.UpdateMailing(ctx, id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Mailing{}, fmt.Errorf("mailingRepository.UpdateMailing: %w", business_errors.ErrMailingNotFound)
		}

		return model.Mailing{}, fmt.Errorf("mailingRepository.UpdateMailing: %w", err)
	}

	return client, nil
}

func (u usecase) DeleteMailing(ctx context.Context, id string) error {
	err := validation.ValidateDeleteMailingRequest(id)
	if err != nil {
		return fmt.Errorf("validateDeleteMailingRequest: %w", err)
	}

	err = u.mailingRepository.DeleteMailing(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("mailingRepository.DeleteMailing: %w", business_errors.ErrMailingNotFound)
		}

		return fmt.Errorf("mailingRepository.DeleteMailing: %w", err)
	}

	return nil
}
