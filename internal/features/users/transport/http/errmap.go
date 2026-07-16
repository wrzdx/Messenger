package users_transport_http

import (
	"errors"
	"messenger/internal/core/domain"
	http_errmap "messenger/internal/core/transport/http/errmap"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func errorMapper(err error) http_response.HTTPError {
	switch {
	case errors.Is(err, domain.ErrInvalidUserProfile):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_user_profile",
			Message:    "invalid user profile",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, domain.ErrAlreadyExists):
		return http_response.HTTPError{
			StatusCode: http.StatusConflict,
			Code:       "user_already_exists",
			Message:    "user already exists",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, domain.ErrNotFound):
		return http_response.HTTPError{
			StatusCode: http.StatusNotFound,
			Code:       "user_not_found",
			Message:    "user not found",
		}

	default:
		return http_errmap.Map(err)
	}
}
