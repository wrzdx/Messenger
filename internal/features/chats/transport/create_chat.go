package chats_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *ChatsHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)
	claims := core_context.ClaimsRequired(ctx)

	var request CreateChatRequest
	if err := http_request.DecodeAndValidateRequest(r, request); err != nil {
		sender.Error(err)
		return
	}

	chat, err := h.chatsService.CreateChat(
		ctx,
		claims.UserID,
		request.Type,
		request.Title,
		request.ParticipantIDs,
	)
	if err != nil {
		sender.Error(err)
		return
	}

	response := CreateChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	}

	sender.OK(http.StatusCreated, response)
}

type CreateChatRequest struct {
	Type           string      `json:"type" validate:"required"`
	Title          *string     `json:"title"`
	ParticipantIDs []uuid.UUID `json:"participant_ids"`
}

type CreateChatResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
