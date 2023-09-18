package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/mailing"
)

type mailingHandler struct {
	mailingUsecase desc.Usecase
}

func (u mailingHandler) CreateMailingHandler(ctx echo.Context) error {
	var req model.CreateMailingRequest

	err := ctx.Bind(&req)
	if err != nil {
		return fmt.Errorf("CreateMailingHandler.Bind: %w", err)
	}

	mailing, err := u.mailingUsecase.CreateMailing(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("mailingUsecase.CreateMailing: %w", err)
	}

	return ctx.JSON(http.StatusCreated, mailing)
}

func (u mailingHandler) UpdateMailingHandler(ctx echo.Context) error {
	var req model.UpdateMailingRequest

	err := ctx.Bind(&req)
	if err != nil {
		return fmt.Errorf("UpdateMailingHandler.Bind: %w", err)
	}

	mailing, err := u.mailingUsecase.UpdateMailing(context.TODO(), ctx.Param("mailingID"), req)
	if err != nil {
		return fmt.Errorf("mailingUsecase.UpdateMailing: %w", err)
	}

	return ctx.JSON(http.StatusOK, mailing)
}

func (u mailingHandler) DeleteMailingHandler(ctx echo.Context) error {
	err := u.mailingUsecase.DeleteMailing(context.TODO(), ctx.Param("mailingID"))
	if err != nil {
		return fmt.Errorf("mailingUsecase.DeleteMailing: %w", err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func NewMailingHandler(e *echo.Echo, mailingUsecase desc.Usecase) mailingHandler {
	handler := mailingHandler{mailingUsecase: mailingUsecase}
	api := e.Group("api/v1")

	createMailingPrefix := "/mailings/create"
	updateMailingPrefix := "/mailings/update/:mailingID"
	deleteMailingPrefix := "/mailings/delete/:mailingID"

	createMailingUrl := api.Group(createMailingPrefix)
	updateMailingUrl := api.Group(updateMailingPrefix)
	deleteMailingUrl := api.Group(deleteMailingPrefix)

	createMailingUrl.POST("", handler.CreateMailingHandler)
	updateMailingUrl.PUT("", handler.UpdateMailingHandler)
	deleteMailingUrl.DELETE("", handler.DeleteMailingHandler)

	return handler
}
