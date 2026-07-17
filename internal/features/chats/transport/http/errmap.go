package chats_transport_http

import (
	"errors"
	"messenger/internal/core/domain"
	http_errmap "messenger/internal/core/transport/http/errmap"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func errorMapper(err error) http_response.HTTPError {
	switch {
	case errors.Is(err, domain.ErrInvalidDirectChat):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_direct_chat",
			Message:    "invalid direct chat",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, domain.ErrNotFound):
		return http_response.HTTPError{
			StatusCode: http.StatusNotFound,
			Code:       "peer_not_found",
			Message:    "peer not found",
			Fields:     http_errmap.FieldsFrom(err),
		}
	default:
		return http_errmap.Map(err)
	}
}
