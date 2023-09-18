package middleware

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/http_utils"
)

type jsonError struct {
	err error `json:"err"`
}

func (j jsonError) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.err.Error())
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestId := uuid.NewString()
		log.Info("Incoming request: ", ctx.Request().URL, ", ip: ", ctx.RealIP(), ", method: ", ctx.Request().Method, ", request_id: ", requestId)

		if err := next(ctx); err != nil {
			statusCode := http_utils.StatusCode(err)
			log.Error("HTTP code: ", statusCode, ", Error: ", err, ", request_id: ", requestId)

			return ctx.JSON(statusCode, jsonError{err: err})
		}

		log.Info("HTTP code: ", ctx.Response().Status, ", request_id: ", requestId)
		return nil
	}
}
