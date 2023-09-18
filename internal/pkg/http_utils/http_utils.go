package http_utils

import (
	"errors"
	"net/http"

	"gitlab.com/fr5270937/notifications_service/internal/pkg/business_errors"
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, business_errors.ErrClientNotFound):
		return http.StatusNotFound
	case errors.Is(err, business_errors.ErrInvalidID):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrNumberFirstDigitDontSeven):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrNumberLenMoreEleven):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrNumberMustConsistOfDigits):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrOperatorCodeMustConsistOfDigits):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrOperatorCodeNotInInterval):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrTagLenMoreHundred):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrTextLenMoreThousand):
		return http.StatusBadRequest
	case errors.Is(err, business_errors.ErrTagLenMoreHundred):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
