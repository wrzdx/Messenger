package chats_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ChatsHandler struct {
	chatsService ChatsService
}

type ChatsService interface {
	CreateChat(
		ctx context.Context,
		userID uuid.UUID,
		chatType string,
		name *string,
		ParticipantIDs []uuid.UUID,
	) (domain.Chat, error)
}

func NewChatsHandler(chatsService ChatsService) *ChatsHandler {
	return &ChatsHandler{
		chatsService: chatsService,
	}
}

func (h *ChatsHandler) Router() chi.Router {
	router := chi.NewRouter()
	return router
}
