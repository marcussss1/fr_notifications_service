package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/labstack/gommon/log"
	"gitlab.com/fr5270937/notifications_service/internal/config"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/business_errors"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/metrics"
	"gitlab.com/fr5270937/notifications_service/internal/services/client"
	"gitlab.com/fr5270937/notifications_service/internal/services/mailing"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/message"
)

type usecase struct {
	metrics           *metrics.Metrics
	clientRepository  client.Repository
	messageRepository desc.Repository
	mailingRepository mailing.Repository
}

func NewMessagesUsecase(metrics *metrics.Metrics, clientRepository client.Repository, messageRepository desc.Repository, mailingRepository mailing.Repository) desc.Usecase {
	return usecase{metrics: metrics, clientRepository: clientRepository, messageRepository: messageRepository, mailingRepository: mailingRepository}
}

func (u usecase) StartMailing(ctx context.Context, mailing model.Mailing) error {
	var clients []model.Client

	// смотрим в фильтре тег или код оператора
	mobileOperatorCode, err := strconv.Atoi(mailing.Filter)
	if err != nil {
		clients, err = u.clientRepository.GetClientsFromTag(ctx, mailing.Filter)
		if err != nil {
			return fmt.Errorf("clientRepository.GetMailingClientsFromTag: %w", err)
		}
	} else {
		clients, err = u.clientRepository.GetClientsFromMobileOperatorCode(ctx, int32(mobileOperatorCode))
		if err != nil {
			return fmt.Errorf("clientRepository.GetMailingClientsFromMobileOperatorCode: %w", err)
		}
	}

	for _, client := range clients {
		go func() {
			message, errIn := u.messageRepository.CreateMessage(ctx, model.Message{
				Status:           model.MESSAGE_STATUS_CREATED,
				MailingID:        mailing.ID,
				ReceiverClientID: client.ID,
			})
			if errIn != nil {
				log.Error(fmt.Errorf("messageRepository.CreateMessage: %w", errIn))
				return
			}
			u.metrics.CountMessages.Inc()

			errIn = u.sendMessageToThirdService(message, mailing, client.PhoneNumber)
			if errIn != nil {
				log.Error(fmt.Errorf("sendMessage: %w", errIn))
			}
		}()
	}

	return nil
}

func (u usecase) GetStatisticMessages(ctx context.Context) ([]model.CommonStatisticMessages, error) {
	mailings, err := u.mailingRepository.GetAllMailings(ctx)
	if err != nil {
		return nil, fmt.Errorf("mailingRepository.GetAllMailings: %w", err)
	}

	var resp []model.CommonStatisticMessages
	for _, mailing := range mailings {
		messages, errIn := u.messageRepository.GetMessagesByMailingID(ctx, mailing.ID)
		if errIn != nil {
			log.Error(fmt.Errorf("messageRepository.GetMessagesByMailingID: %w", errIn))
			continue
		}

		resp = append(resp, model.CommonStatisticMessages{
			Mailing:  mailing,
			Messages: messages,
		})
	}

	return resp, nil
}

func (u usecase) GetDetailStatisticMessagesByMailing(ctx context.Context, id string) ([]model.Message, error) {
	messages, err := u.messageRepository.GetMessagesByMailingID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("messageRepository.GetMessagesByMailingID: %w", err)
	}

	return messages, err
}

func (u usecase) HandlingActiveMessages(ctx context.Context) error {
	mailings, err := u.mailingRepository.GetActiveMailings(ctx)
	if err != nil {
		return fmt.Errorf("mailingRepository.GetActiveMailings: %w", err)
	}

	for _, mailing := range mailings {
		messages, errIn := u.messageRepository.GetStatusMessagesByMailingIDAndStatus(ctx, mailing.ID, model.MESSAGE_STATUS_CREATED)
		if errIn != nil {
			log.Error(fmt.Errorf("messageRepository.GetStatusMessagesByMailingIDAndStatus: %w", errIn))
			continue
		}

		for _, message := range messages {
			go func() {
				err = u.sendMessageFromClient(ctx, message, mailing)
				if err != nil {
					log.Error(fmt.Errorf("sendMessageFromClient: %w", err))
					return
				}

				u.metrics.CountFailedMessages.Add(-1)
			}()
		}
	}

	return nil
}

func (u usecase) sendMessageToThirdService(message model.Message, mailing model.Mailing, phoneNumber string) error {
	phoneNumberNum, err := strconv.Atoi(phoneNumber)
	if err != nil {
		return fmt.Errorf("strconv.Atoi(phoneNumber): %w", err)
	}

	// поход в сторонний сервис
	resp, err := u.sendRetryRequestsToThirdService(message, mailing, phoneNumberNum)
	if err != nil || resp.StatusCode != http.StatusOK {
		// рассылка истекла - удаляем из таблицы отложенных рассылок
		if errors.Is(err, business_errors.ErrMailingTimeIsUp) {
			errIn := u.mailingRepository.DeletePendingMailing(context.TODO(), mailing.ID)
			if errIn != nil {
				return fmt.Errorf("mailingRepository.DeletePendingMailing: %w", errIn)
			}

			return fmt.Errorf("sendRetryRequests: %w", business_errors.ErrMailingTimeIsUp)
		}

		u.metrics.CountFailedMessages.Inc()
		_, errIn := u.messageRepository.UpdateMessage(context.TODO(), message.ID, model.Message{
			ID:     message.ID,
			Status: model.MESSAGE_STATUS_FAILED,
		})
		if errIn != nil {
			return fmt.Errorf("messageRepository.UpdateMessage: %w", errIn)
		}

		return fmt.Errorf("sendRetryRequests: %w", err)
	}

	u.metrics.CountSucceededMessages.Inc()
	_, err = u.messageRepository.UpdateMessage(context.TODO(), message.ID, model.Message{
		ID:     message.ID,
		Status: model.MESSAGE_STATUS_SUCCEEDED,
	})

	return err
}

func (u usecase) sendRetryRequestsToThirdService(message model.Message, mailing model.Mailing, phoneNumber int) (*http.Response, error) {
	req, err := u.createRequestToThirdService(message, mailing, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("createRequestToThirdService: %w", err)
	}

	// три попытки, интервалы между попытками: 1, 2, 4 секунды
	// в зависимости от бизнес-требований логика ретраев может меняться
	client := &http.Client{}
	retryOptions := []retry.Option{
		retry.Attempts(3),
		retry.Delay(time.Second * 1),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("Попытка %d: %v\n", n, err)
		}),
		retry.RetryIf(func(err error) bool {
			if errors.Is(err, business_errors.ErrMailingTimeIsUp) {
				return false
			}

			return true
		}),
	}

	var resp *http.Response
	err = retry.Do(func() error {
		if time.Now().String() > mailing.FinishedAt.String() {
			return business_errors.ErrMailingTimeIsUp
		}

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("client.Do: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()

		return nil
	}, retryOptions...)
	if err != nil {
		return nil, fmt.Errorf("retry.Do: %w", err)
	}

	return resp, nil
}

func (u usecase) createRequestToThirdService(message model.Message, mailing model.Mailing, phoneNumber int) (*http.Request, error) {
	sendMessageData := model.SendMessageRequest{
		ID:          message.ID,
		PhoneNumber: int64(phoneNumber),
		Text:        mailing.Text,
	}
	jsonData, err := json.Marshal(sendMessageData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequest("POST", config.ApiUrl+strconv.Itoa(int(message.ID)), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.JWT)

	return req, nil
}

func (u usecase) sendMessageFromClient(ctx context.Context, message model.Message, mailing model.Mailing) error {
	client, err := u.clientRepository.GetClientFromID(ctx, message.ReceiverClientID)
	if err != nil {
		return fmt.Errorf("clientRepository.GetClientFromID: %w", err)
	}

	err = u.sendMessageToThirdService(message, mailing, client.PhoneNumber)
	if err != nil {
		return fmt.Errorf("sendMessageToThirdService: %w", err)
	}

	return nil
}
