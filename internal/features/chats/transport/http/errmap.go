package chats_transport_http

import (
	"errors"
	"messenger/internal/core/domain"
	http_errmap "messenger/internal/core/transport/http/errmap"
	http_response "messenger/internal/core/transport/http/response"
	chats_service "messenger/internal/features/chats/service"
	"net/http"
)

func errorMapper(err error) http_response.HTTPError {
	switch {
	case errors.Is(err, chats_service.ErrInvalidListChatsQuery):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_list_chats_query",
			Message:    "invalid list chats query",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, domain.ErrInvalidDirectChat):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_direct_chat",
			Message:    "invalid direct chat",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, domain.ErrInvalidGroupChat):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_group_chat",
			Message:    "invalid group chat",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, chats_service.ErrInvalidListGroupParticipantsQuery):
		return http_response.HTTPError{
			StatusCode: http.StatusBadRequest,
			Code:       "invalid_list_group_participants_query",
			Message:    "invalid list group participants query",
			Fields:     http_errmap.FieldsFrom(err),
		}

	case errors.Is(err, domain.ErrNotFound):
		return http_response.HTTPError{
			StatusCode: http.StatusNotFound,
			Code:       "not_found",
			Message:    "not found",
			Fields:     http_errmap.FieldsFrom(err),
		}
	default:
		return http_errmap.Map(err)
	}
}
