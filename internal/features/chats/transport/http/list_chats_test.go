package chats_transport_http

import (
	"encoding/json"
	"errors"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	"messenger/internal/core/logger"
	http_middleware "messenger/internal/core/transport/http/middleware"
	http_response "messenger/internal/core/transport/http/response"
	chats_service "messenger/internal/features/chats/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListChats(t *testing.T) {
	requesterID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	peerID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	activityAt := time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)

	t.Run("returns direct and group chat summaries", func(t *testing.T) {
		direct := newChatsTransportDirectItem(t, requesterID, peerID, activityAt)
		group := newChatsTransportGroupItem(t, requesterID, activityAt.Add(-time.Minute))
		before := &chats_service.ChatCursor{ChatID: uuid.New(), LastActivityAt: activityAt.Add(time.Hour)}
		next := &chats_service.ChatCursor{
			ChatID:         group.Group.Chat.ID,
			LastActivityAt: group.Group.Chat.LastActivityAt,
		}
		service := NewMockChatsService(t)
		service.EXPECT().ListChats(
			mock.Anything,
			requesterID,
			chats_service.ListChatsQuery{Before: before, Limit: 2},
		).Return(chats_service.ChatPage{
			Chats:      []chats_service.ChatItem{direct, group},
			NextCursor: next,
		}, nil)
		router := newChatsTransportRouter(service, requesterID)
		encodedBefore, err := encodeChatCursor(before)
		require.NoError(t, err)
		request := newListChatsRequest(t, "?limit=2&cursor="+url.QueryEscape(*encodedBefore))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		response := decodeListChatsResponse(t, recorder)
		require.Len(t, response.Chats, 2)
		require.NotNil(t, response.Chats[0].Direct)
		require.Nil(t, response.Chats[0].Group)
		require.Equal(t, peerID, response.Chats[0].Direct.Peer.ID)
		require.Equal(t, direct.Direct.PeerProfile.Username(), response.Chats[0].Direct.Peer.Username)
		require.Nil(t, response.Chats[0].LastMessage)
		require.NotNil(t, response.Chats[1].Group)
		require.Equal(t, group.Group.Title, response.Chats[1].Group.Title)
		require.NotNil(t, response.Chats[1].LastMessage)
		require.Equal(
			t,
			group.LastMessage.SenderProfile.Username(),
			response.Chats[1].LastMessage.SenderUsername,
		)
		require.NotNil(t, response.NextCursor)
		decodedNext, err := decodeChatCursor(*response.NextCursor)
		require.NoError(t, err)
		require.Equal(t, next, decodedNext)
	})

	t.Run("returns empty final page", func(t *testing.T) {
		service := NewMockChatsService(t)
		service.EXPECT().ListChats(
			mock.Anything,
			requesterID,
			chats_service.ListChatsQuery{},
		).Return(chats_service.ChatPage{Chats: []chats_service.ChatItem{}}, nil)
		router := newChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, newListChatsRequest(t, ""))

		require.Equal(t, http.StatusOK, recorder.Code)
		response := decodeListChatsResponse(t, recorder)
		require.NotNil(t, response.Chats)
		require.Empty(t, response.Chats)
		require.Nil(t, response.NextCursor)
	})

	testCases := []struct {
		name  string
		query string
	}{
		{name: "malformed cursor", query: "?cursor=%25%25%25"},
		{name: "malformed limit", query: "?limit=many"},
	}
	for _, testCase := range testCases {
		t.Run("rejects "+testCase.name, func(t *testing.T) {
			router := newChatsTransportRouter(NewMockChatsService(t), requesterID)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, newListChatsRequest(t, testCase.query))

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Equal(t, http_response.APIErrorDetail{
				Code:    "invalid_request",
				Message: "invalid request",
			}, decodeChatsTransportError(t, recorder))
		})
	}

	t.Run("maps invalid service query", func(t *testing.T) {
		service := NewMockChatsService(t)
		service.EXPECT().ListChats(
			mock.Anything,
			requesterID,
			chats_service.ListChatsQuery{Limit: -1},
		).Return(chats_service.ChatPage{}, domain.DetailedError{
			Err: chats_service.ErrInvalidListChatsQuery,
			Details: map[string]string{
				"limit": "limit must be between 1 and 100",
			},
		})
		router := newChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, newListChatsRequest(t, "?limit=-1"))

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "invalid_list_chats_query",
			Message: "invalid list chats query",
			Fields: map[string]string{
				"limit": "limit must be between 1 and 100",
			},
		}, decodeChatsTransportError(t, recorder))
	})

	t.Run("does not expose unexpected service error", func(t *testing.T) {
		serviceErr := errors.New("database unavailable")
		service := NewMockChatsService(t)
		service.EXPECT().ListChats(
			mock.Anything,
			requesterID,
			chats_service.ListChatsQuery{},
		).Return(chats_service.ChatPage{}, serviceErr)
		router := newChatsTransportRouter(service, requesterID)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, newListChatsRequest(t, ""))

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.Equal(t, http_response.APIErrorDetail{
			Code:    "internal_error",
			Message: "internal server error",
		}, decodeChatsTransportError(t, recorder))
		require.NotContains(t, recorder.Body.String(), serviceErr.Error())
	})
}

func newChatsTransportRouter(service ChatsService, currentUserID uuid.UUID) http.Handler {
	authMiddleware := http_middleware.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := core_context.WithClaims(
				r.Context(),
				core_context.ContextClaims{UserID: currentUserID},
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	router := chi.NewRouter()
	router.Mount("/chats", NewChatsHandler(service).Router(authMiddleware))
	return router
}

func newListChatsRequest(t *testing.T, query string) *http.Request {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/chats"+query, nil)
	ctx := logger.WithLogger(request.Context(), logger.NewTestLogger())
	return request.WithContext(ctx)
}

func decodeListChatsResponse(t *testing.T, recorder *httptest.ResponseRecorder) ListChatsResponse {
	t.Helper()
	var body struct {
		Data ListChatsResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	return body.Data
}

func newChatsTransportDirectItem(
	t *testing.T,
	requesterID uuid.UUID,
	peerID uuid.UUID,
	activityAt time.Time,
) chats_service.ChatItem {
	t.Helper()
	direct, err := domain.NewDirectChat(uuid.New(), requesterID, peerID, activityAt)
	require.NoError(t, err)
	profile, err := domain.NewUserProfile("Peer_123", "Peer", nil, nil)
	require.NoError(t, err)
	item := chats_service.ChatItem{Direct: &chats_service.DirectChatItem{
		Chat:        direct,
		PeerID:      peerID,
		PeerProfile: profile,
	}}
	require.NoError(t, item.Validate())
	return item
}

func newChatsTransportGroupItem(
	t *testing.T,
	senderID uuid.UUID,
	activityAt time.Time,
) chats_service.ChatItem {
	t.Helper()
	group, err := domain.NewGroupChat(uuid.New(), "Backend", activityAt.Add(-time.Hour))
	require.NoError(t, err)
	message, err := domain.NewMessage(
		uuid.New(), uuid.New(), group.Chat.ID, senderID, "group message", activityAt,
	)
	require.NoError(t, err)
	group.Chat.LastMessageID = &message.ID
	group.Chat.LastActivityAt = message.CreatedAt
	senderProfile, err := domain.NewUserProfile("Author_123", "Author", nil, nil)
	require.NoError(t, err)
	item := chats_service.ChatItem{
		Group: &group,
		LastMessage: &chats_service.LastMessageItem{
			Message:       message,
			SenderProfile: senderProfile,
		},
	}
	require.NoError(t, item.Validate())
	return item
}
