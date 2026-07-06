package users_transport_http

import (
	"errors"
	"messenger/internal/core/domain"
	"net/http"
)

var (
	ErrMissingClaims   = errors.New("missing claims")
	ErrInvalidArgument = errors.New("invalid argument")
)

func mapDomainErrorToStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidUsername),
		errors.Is(err, domain.ErrInvalidFirstName),
		errors.Is(err, domain.ErrInvalidLastName),
		errors.Is(err, domain.ErrInvalidBio):
		return http.StatusBadRequest

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict

	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
