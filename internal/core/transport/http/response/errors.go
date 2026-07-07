package core_http_response

import (
	"errors"
	"messenger/internal/core/domain"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"   example:"Error description"`
}

type Error struct {
	Error   error
	Status  int
	Message string
}

var (
	ErrMissingClaims   = errors.New("missing claims")
	ErrInvalidArgument = errors.New("invalid argument")
)

var errorMap = []struct {
	err    error
	status int
}{
	{ErrMissingClaims, http.StatusUnauthorized},
	{ErrInvalidArgument, http.StatusBadRequest},
	{domain.ErrInvalidUsername, http.StatusBadRequest},
	{domain.ErrInvalidFirstName, http.StatusBadRequest},
	{domain.ErrInvalidLastName, http.StatusBadRequest},
	{domain.ErrInvalidBio, http.StatusBadRequest},
	{domain.ErrNegativeOffset, http.StatusBadRequest},
	{domain.ErrNegativeLimit, http.StatusBadRequest},
	{domain.ErrNullUsername, http.StatusBadRequest},
	{domain.ErrNullFirstname, http.StatusBadRequest},
	{domain.ErrInvalidPassword, http.StatusBadRequest},

	{domain.ErrUserAlreadyExists, http.StatusConflict},
	{domain.ErrInvalidCredentials, http.StatusUnauthorized},
	{domain.ErrUserNotFound, http.StatusNotFound},
}

func MapError(err error) Error {
	for _, e := range errorMap {
		if errors.Is(err, e.err) {
			return Error{
				Error:   e.err,
				Status:  e.status,
				Message: e.err.Error(),
			}
		}
	}
	return Error{
		Error:   err,
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}
}
