package core_http_response

import (
	"errors"
	"messenger/internal/core/domain"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"   example:"Error description"`
}

var (
	ErrMissingClaims   = errors.New("missing claims")
	ErrInvalidArgument = errors.New("invalid argument")
)

func MapDomainErrorToStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidUsername),
		errors.Is(err, domain.ErrInvalidFirstName),
		errors.Is(err, domain.ErrInvalidLastName),
		errors.Is(err, domain.ErrInvalidBio),
		errors.Is(err, domain.ErrNegativeOffset),
		errors.Is(err, domain.ErrNegativeLimit):
		return http.StatusBadRequest

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict

	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized

	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
