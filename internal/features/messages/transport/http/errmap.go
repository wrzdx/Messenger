package messages_transport_http

import (
	"errors"
	"messenger/internal/core/domain"
	http_errmap "messenger/internal/core/transport/http/errmap"
	http_response "messenger/internal/core/transport/http/response"
	messages_service "messenger/internal/features/messages/service"
	"net/http"
)

func errorMapper(err error) http_response.HTTPError {
	switch {
	case errors.Is(err, domain.ErrInvalidMessage):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_message",
			Message:    "invalid message",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, messages_service.ErrMessageConflict):
		return http_response.HTTPError{
			StatusCode: http.StatusConflict,
			Code:       "message_conflict",
			Message:    "message conflict",
		}

	case errors.Is(err, messages_service.ErrMessageTargetUnavailable):
		return http_response.HTTPError{
			StatusCode: http.StatusConflict,
			Code:       "message_target_unavailable",
			Message:    "message target unavailable",
		}

	case errors.Is(err, domain.ErrNotFound):
		return http_response.HTTPError{
			StatusCode: http.StatusNotFound,
			Code:       "chat_not_found",
			Message:    "chat not found",
			Fields:     http_errmap.FieldsFrom(err),
		}
	default:
		return http_errmap.Map(err)
	}
}
