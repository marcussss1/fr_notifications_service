package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/client"
)

type clientHandler struct {
	clientUsecase desc.Usecase
}

func (u clientHandler) CreateClientHandler(ctx echo.Context) error {
	var req model.CreateClientRequest

	err := ctx.Bind(&req)
	if err != nil {
		return fmt.Errorf("CreateClientHandler.Bind: %w", err)
	}

	client, err := u.clientUsecase.CreateClient(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("clientUsecase.CreateClient: %w", err)
	}

	return ctx.JSON(http.StatusCreated, client)
}

func (u clientHandler) UpdateClientHandler(ctx echo.Context) error {
	var req model.UpdateClientRequest

	err := ctx.Bind(&req)
	if err != nil {
		return fmt.Errorf("UpdateClientHandler.Bind: %w", err)
	}

	client, err := u.clientUsecase.UpdateClient(context.TODO(), ctx.Param("clientID"), req)
	if err != nil {
		return fmt.Errorf("clientUsecase.UpdateClient: %w", err)
	}

	return ctx.JSON(http.StatusOK, client)
}

func (u clientHandler) DeleteClientHandler(ctx echo.Context) error {
	err := u.clientUsecase.DeleteClient(context.TODO(), ctx.Param("clientID"))
	if err != nil {
		return fmt.Errorf("clientUsecase.DeleteClient: %w", err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func NewClientHandler(e *echo.Echo, clientUsecase desc.Usecase) clientHandler {
	handler := clientHandler{clientUsecase: clientUsecase}
	api := e.Group("api/v1")

	createClientPrefix := "/clients/create"
	updateClientPrefix := "/clients/update/:clientID"
	deleteClientPrefix := "/clients/delete/:clientID"

	createClientUrl := api.Group(createClientPrefix)
	updateClientUrl := api.Group(updateClientPrefix)
	deleteClientUrl := api.Group(deleteClientPrefix)

	createClientUrl.POST("", handler.CreateClientHandler)
	updateClientUrl.PUT("", handler.UpdateClientHandler)
	deleteClientUrl.DELETE("", handler.DeleteClientHandler)

	return handler
}
