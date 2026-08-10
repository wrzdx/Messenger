package chats_transport_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"messenger/internal/core/domain"
	"messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	chats_service "messenger/internal/features/chats/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddGroupParticipants(t *testing.T) {
	requesterID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	memberID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	unavailableID := uuid.MustParse("00000000-0000-0000-0000-000000000004")

	t.Run("returns positional participant results", func(t *testing.T) {
		participantIDs := []uuid.UUID{memberID, unavailableID, memberID}
		command := chats_service.AddGroupParticipantsCommand{
			GroupID:        groupID,
			RequesterID:    requesterID,
			ParticipantIDs: participantIDs,
		}
		result := []chats_service.AddGroupParticipantResult{
			{UserID: memberID, Status: "added"},
			{UserID: unavailableID, Status: "unavailable"},
			{UserID: memberID, Status: "already_member"},
		}
		service := NewMockChatsService(t)
		service.EXPECT().
			AddGroupParticipants(mock.Anything, command).
			Return(result, nil)
		router := newListChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": participantIDs,
			}),
		)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		require.Equal(t, AddGroupParticipantsResponse{
			{UserID: memberID, Status: "added"},
			{UserID: unavailableID, Status: "unavailable"},
			{UserID: memberID, Status: "already_member"},
		}, decodeAddGroupParticipantsResponse(t, recorder))
	})

	t.Run("allows empty participant list", func(t *testing.T) {
		service := NewMockChatsService(t)
		service.EXPECT().AddGroupParticipants(
			mock.Anything,
			chats_service.AddGroupParticipantsCommand{
				GroupID:        groupID,
				RequesterID:    requesterID,
				ParticipantIDs: []uuid.UUID{},
			},
		).Return([]chats_service.AddGroupParticipantResult{}, nil)
		router := newListChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": []uuid.UUID{},
			}),
		)

		require.Equal(t, http.StatusOK, recorder.Code)
		response := decodeAddGroupParticipantsResponse(t, recorder)
		require.NotNil(t, response)
		require.Empty(t, response)
	})

	t.Run("rejects malformed chat id without calling service", func(t *testing.T) {
		router := newListChatsTransportRouter(NewMockChatsService(t), requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, "not-a-uuid", map[string]any{
				"participant_ids": []uuid.UUID{memberID},
			}),
		)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "invalid_request",
			Message: "invalid request",
			Fields: map[string]string{
				"chat_id": "invalid uuid",
			},
		}, decodeChatsTransportError(t, recorder))
	})

	t.Run("rejects malformed participant id without calling service", func(t *testing.T) {
		router := newListChatsTransportRouter(NewMockChatsService(t), requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": []string{memberID.String(), "not-a-uuid"},
			}),
		)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "invalid_request",
			Message: "invalid request",
			Fields: map[string]string{
				"participant_ids[1]": "invalid uuid",
			},
		}, decodeChatsTransportError(t, recorder))
	})

	t.Run("maps invalid service command", func(t *testing.T) {
		command := chats_service.AddGroupParticipantsCommand{
			GroupID:        groupID,
			RequesterID:    requesterID,
			ParticipantIDs: []uuid.UUID{memberID},
		}
		service := NewMockChatsService(t)
		service.EXPECT().
			AddGroupParticipants(mock.Anything, command).
			Return(nil, chats_service.ErrInvalidAddGroupParticipantsQuery)
		router := newListChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": []uuid.UUID{memberID},
			}),
		)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "invalid_add_group_participants_query",
			Message: "invalid add group participants query",
		}, decodeChatsTransportError(t, recorder))
	})

	t.Run("maps insufficient rights", func(t *testing.T) {
		command := chats_service.AddGroupParticipantsCommand{
			GroupID:        groupID,
			RequesterID:    requesterID,
			ParticipantIDs: []uuid.UUID{memberID},
		}
		service := NewMockChatsService(t)
		service.EXPECT().
			AddGroupParticipants(mock.Anything, command).
			Return(nil, chats_service.ErrNotEnoughRights)
		router := newListChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": []uuid.UUID{memberID},
			}),
		)

		require.Equal(t, http.StatusForbidden, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "not_enough_rights",
			Message: "not enough rights",
		}, decodeChatsTransportError(t, recorder))
	})

	t.Run("maps missing group or requester", func(t *testing.T) {
		command := chats_service.AddGroupParticipantsCommand{
			GroupID:        groupID,
			RequesterID:    requesterID,
			ParticipantIDs: []uuid.UUID{memberID},
		}
		service := NewMockChatsService(t)
		service.EXPECT().
			AddGroupParticipants(mock.Anything, command).
			Return(nil, domain.ErrNotFound)
		router := newListChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": []uuid.UUID{memberID},
			}),
		)

		require.Equal(t, http.StatusNotFound, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "not_found",
			Message: "not found",
		}, decodeChatsTransportError(t, recorder))
	})

	t.Run("does not expose unexpected service error", func(t *testing.T) {
		serviceErr := errors.New("database unavailable")
		command := chats_service.AddGroupParticipantsCommand{
			GroupID:        groupID,
			RequesterID:    requesterID,
			ParticipantIDs: []uuid.UUID{memberID},
		}
		service := NewMockChatsService(t)
		service.EXPECT().
			AddGroupParticipants(mock.Anything, command).
			Return(nil, serviceErr)
		router := newListChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			newAddGroupParticipantsHTTPRequest(t, groupID.String(), map[string]any{
				"participant_ids": []uuid.UUID{memberID},
			}),
		)

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "internal_error",
			Message: "internal server error",
		}, decodeChatsTransportError(t, recorder))
		require.NotContains(t, recorder.Body.String(), serviceErr.Error())
	})
}

func newAddGroupParticipantsHTTPRequest(
	t *testing.T,
	groupID string,
	body any,
) *http.Request {
	t.Helper()
	payload, err := json.Marshal(body)
	require.NoError(t, err)
	request := httptest.NewRequest(
		http.MethodPost,
		"/chats/groups/"+groupID+"/participants",
		bytes.NewReader(payload),
	)
	ctx := logger.WithLogger(request.Context(), logger.NewTestLogger())
	return request.WithContext(ctx)
}

func decodeAddGroupParticipantsResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
) AddGroupParticipantsResponse {
	t.Helper()
	var body struct {
		Data AddGroupParticipantsResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	return body.Data
}
