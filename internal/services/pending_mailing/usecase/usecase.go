package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	"gitlab.com/fr5270937/notifications_service/internal/services/mailing"
	"gitlab.com/fr5270937/notifications_service/internal/services/message"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/pending_mailing"
)

type usecase struct {
	messageUsecase    message.Usecase
	mailingRepository mailing.Repository
}

func NewPendingMailingUsecase(messageUsecase message.Usecase, mailingRepository mailing.Repository) desc.Usecase {
	return usecase{messageUsecase: messageUsecase, mailingRepository: mailingRepository}
}

func (u usecase) PendingMailingWatchdog() {
	for {
		mailings, err := u.mailingRepository.GetMailingsToSending(context.TODO())
		if err != nil {
			log.Error(fmt.Errorf("mailingRepository.GetMailingsToSending: %w", err))
			time.Sleep(1 * time.Second)
			continue
		}

		for _, mailing := range mailings {
			go func() {
				errIn := u.updatePendingMailingStatus(mailing.ID, model.PENDING_MAILING_STATUS_IN_WORK)
				if errIn != nil {
					log.Error(fmt.Errorf("updatePendingMailingStatus: %w", errIn))
					return
				}

				errIn = u.messageUsecase.StartMailing(context.TODO(), mailing)
				if errIn != nil {
					log.Error(fmt.Errorf("messageUsecase.StartMailing: %w", errIn))
					return
				}

				// если рассылка прошла успешно - удаляем ее из таблицы отложенных рассылок
				errIn = u.mailingRepository.DeletePendingMailing(context.TODO(), mailing.ID)
				if errIn != nil {
					log.Error(fmt.Errorf("mailingRepository.DeletePendingMailing: %w", errIn))
				}
			}()
		}

		time.Sleep(1 * time.Second)
	}
}

func (u usecase) updatePendingMailingStatus(mailingID string, status int) error {
	_, err := u.mailingRepository.UpdatePendingMailing(context.TODO(), mailingID, model.UpdatePendingMailing{
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("mailingRepository.UpdatePendingMailing: %w", err)
	}

	return nil
}
