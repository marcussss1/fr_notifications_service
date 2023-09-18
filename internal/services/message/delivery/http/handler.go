package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/message"
)

type messagesHandler struct {
	messagesUsecase desc.Usecase
}

func (u messagesHandler) GetStatisticMessagesHandler(ctx echo.Context) error {
	statistic, err := u.messagesUsecase.GetStatisticMessages(context.TODO())
	if err != nil {
		return fmt.Errorf("messagesUsecase.GetStatisticMessages: %w", err)
	}

	return ctx.JSON(http.StatusOK, statistic)
}

func (u messagesHandler) GetDetailStatisticMessagesByMailingHandler(ctx echo.Context) error {
	statistic, err := u.messagesUsecase.GetDetailStatisticMessagesByMailing(context.TODO(), ctx.Param("mailingID"))
	if err != nil {
		return fmt.Errorf("messagesUsecase.GetDetailStatisticMessagesByMailing: %w", err)
	}

	return ctx.JSON(http.StatusOK, statistic)
}

func (u messagesHandler) HandlingActiveMessagesHandler(ctx echo.Context) error {
	err := u.messagesUsecase.HandlingActiveMessages(context.TODO())
	if err != nil {
		return fmt.Errorf("messagesUsecase.HandlingActiveMessages: %w", err)
	}

	return ctx.NoContent(http.StatusOK)
}

func NewMessagesHandler(e *echo.Echo, messagesUsecase desc.Usecase) messagesHandler {
	handler := messagesHandler{messagesUsecase: messagesUsecase}
	api := e.Group("api/v1")

	getStatisticMessagesPrefix := "/messages"
	getDetailStatisticFromMailingIDPrefix := "/messages/:mailingID"
	handlingActiveMessagesPrefix := "/messages/active"

	getStatisticMessagesUrl := api.Group(getStatisticMessagesPrefix)
	getDetailStatisticFromMailingIDUrl := api.Group(getDetailStatisticFromMailingIDPrefix)
	handlingActiveMessagesUrl := api.Group(handlingActiveMessagesPrefix)

	getStatisticMessagesUrl.GET("", handler.GetStatisticMessagesHandler)
	getDetailStatisticFromMailingIDUrl.GET("", handler.GetDetailStatisticMessagesByMailingHandler)
	handlingActiveMessagesUrl.PUT("", handler.HandlingActiveMessagesHandler)

	return handler
}
