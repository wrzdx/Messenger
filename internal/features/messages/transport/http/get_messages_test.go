package messages_transport_http

import (
	"encoding/json"
	"errors"
	"messenger/internal/core/domain"
	"messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	messages_service "messenger/internal/features/messages/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMessages(t *testing.T) {
	currentUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	chatID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	cursor := &messages_service.MessageCursor{
		MessageID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		CreatedAt: time.Date(2026, time.July, 19, 12, 0, 0, 0, time.UTC),
	}

	t.Run("returns message page", func(t *testing.T) {
		message := newMessagesTransportMessage(t, currentUserID, messages_service.SendMessageCommand{
			ChatID:          chatID,
			ClientMessageID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Content:         "Hello",
		})
		nextCursor := &messages_service.MessageCursor{
			MessageID: message.ID,
			CreatedAt: message.CreatedAt,
		}
		service := NewMockMessagesService(t)
		service.EXPECT().GetMessages(
			mock.Anything,
			currentUserID,
			messages_service.GetMessagesQuery{
				ChatID: chatID,
				Before: cursor,
				Limit:  25,
			},
		).Return(messages_service.MessagePage{
			Messages:   []domain.Message{message},
			NextCursor: nextCursor,
		}, nil)
		router := newMessagesTransportRouter(service, currentUserID)
		encodedCursor, err := encodeMessageCursor(cursor)
		require.NoError(t, err)
		request := newGetMessagesRequest(
			t,
			chatID.String(),
			"?limit=25&cursor="+url.QueryEscape(*encodedCursor),
		)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		response := decodeGetMessagesResponse(t, recorder)
		require.Len(t, response.Messages, 1)
		requireMessagesTransportResponseEqual(
			t,
			message,
			SendMessageResponse(response.Messages[0]),
		)
		require.NotNil(t, response.NextCursor)
		decodedNextCursor, err := decodeMessageCursor(*response.NextCursor)
		require.NoError(t, err)
		require.Equal(t, nextCursor, decodedNextCursor)
	})

	t.Run("returns empty final page with null cursor", func(t *testing.T) {
		service := NewMockMessagesService(t)
		service.EXPECT().GetMessages(
			mock.Anything,
			currentUserID,
			messages_service.GetMessagesQuery{ChatID: chatID},
		).Return(messages_service.MessagePage{
			Messages: []domain.Message{},
		}, nil)
		router := newMessagesTransportRouter(service, currentUserID)
		request := newGetMessagesRequest(t, chatID.String(), "")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		response := decodeGetMessagesResponse(t, recorder)
		require.Empty(t, response.Messages)
		require.NotNil(t, response.Messages)
		require.Nil(t, response.NextCursor)
	})

	testCases := []struct {
		name     string
		chatID   string
		query    string
		response http_response.APIErrorDetail
	}{
		{
			name:   "rejects malformed chat id",
			chatID: "not-a-uuid",
			response: http_response.APIErrorDetail{
				Code:    "invalid_request",
				Message: "invalid request",
				Fields: map[string]string{
					"chat_id": "invalid uuid",
				},
			},
		},
		{
			name:   "rejects malformed limit",
			chatID: chatID.String(),
			query:  "?limit=many",
			response: http_response.APIErrorDetail{
				Code:    "invalid_request",
				Message: "invalid request",
				Fields: map[string]string{
					"limit": "invalid limit query param",
				},
			},
		},
		{
			name:   "rejects malformed cursor",
			chatID: chatID.String(),
			query:  "?cursor=%25%25%25",
			response: http_response.APIErrorDetail{
				Code:    "invalid_request",
				Message: "invalid request",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			router := newMessagesTransportRouter(NewMockMessagesService(t), currentUserID)
			request := newGetMessagesRequest(t, testCase.chatID, testCase.query)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Equal(t, testCase.response, decodeMessagesTransportError(t, recorder))
		})
	}

	t.Run("maps invalid service query", func(t *testing.T) {
		service := NewMockMessagesService(t)
		service.EXPECT().GetMessages(
			mock.Anything,
			currentUserID,
			messages_service.GetMessagesQuery{ChatID: chatID, Limit: -1},
		).Return(messages_service.MessagePage{}, domain.DetailedError{
			Err: messages_service.ErrInvalidGetMessagesQuery,
			Details: map[string]string{
				"limit": "limit must be between 1 and 100",
			},
		})
		router := newMessagesTransportRouter(service, currentUserID)
		request := newGetMessagesRequest(t, chatID.String(), "?limit=-1")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "invalid_get_message_query",
			Message: "invalid get message query",
			Fields: map[string]string{
				"limit": "limit must be between 1 and 100",
			},
		}, decodeMessagesTransportError(t, recorder))
	})

	t.Run("returns chat not found", func(t *testing.T) {
		service := NewMockMessagesService(t)
		service.EXPECT().GetMessages(
			mock.Anything,
			currentUserID,
			messages_service.GetMessagesQuery{ChatID: chatID},
		).Return(messages_service.MessagePage{}, domain.ErrNotFound)
		router := newMessagesTransportRouter(service, currentUserID)
		request := newGetMessagesRequest(t, chatID.String(), "")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusNotFound, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "chat_not_found",
			Message: "chat not found",
		}, decodeMessagesTransportError(t, recorder))
	})

	t.Run("does not expose unexpected service error", func(t *testing.T) {
		serviceErr := errors.New("database unavailable")
		service := NewMockMessagesService(t)
		service.EXPECT().GetMessages(
			mock.Anything,
			currentUserID,
			messages_service.GetMessagesQuery{ChatID: chatID},
		).Return(messages_service.MessagePage{}, serviceErr)
		router := newMessagesTransportRouter(service, currentUserID)
		request := newGetMessagesRequest(t, chatID.String(), "")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "internal_error",
			Message: "internal server error",
		}, decodeMessagesTransportError(t, recorder))
		require.NotContains(t, recorder.Body.String(), serviceErr.Error())
	})
}

func newGetMessagesRequest(
	t *testing.T,
	chatID string,
	query string,
) *http.Request {
	t.Helper()
	request := httptest.NewRequest(
		http.MethodGet,
		"/chats/"+chatID+"/messages"+query,
		nil,
	)
	ctx := logger.WithLogger(request.Context(), logger.NewTestLogger())
	return request.WithContext(ctx)
}

func decodeGetMessagesResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
) GetMessagesResponse {
	t.Helper()
	var body struct {
		Data GetMessagesResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	return body.Data
}
