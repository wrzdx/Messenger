package messages_transport_http

import (
	"fmt"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	messages_service "messenger/internal/features/messages/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *MessagesHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	var request SendMessageRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := uuid.Parse(chatIDStr)

	if err != nil {
		sender.Error(fmt.Errorf("invalid chat id: %w", http_request.ErrInvalidRequest))
		return
	}

	message, isCreated, err := h.messagesService.SendMessage(
		ctx,
		claims.UserID,
		messages_service.SendMessageCommand{
			ChatID:          chatID,
			ClientMessageID: request.ClientMessageID,
			Content:         request.Content,
		},
	)
	if err != nil {
		sender.Error(err)
		return
	}
	response := SendMessageResponse{
		ID:              message.ID,
		ClientMessageID: message.ClientMessageID,
		ChatID:          message.ChatID,
		SenderID:        message.SenderID,
		Content:         message.Content,
		CreatedAt:       message.CreatedAt,
		UpdatedAt:       message.UpdatedAt,
	}
	status := http.StatusOK
	if isCreated {
		status = http.StatusCreated
	}

	sender.OK(status, response)
}

type SendMessageRequest struct {
	ClientMessageID uuid.UUID `json:"client_message_id" validate:"required"`
	Content         string    `json:"content" validate:"required"`
}

type SendMessageResponse struct {
	ID              uuid.UUID  `json:"id"`
	ClientMessageID uuid.UUID  `json:"client_message_id"`
	ChatID          uuid.UUID  `json:"chat_id"`
	SenderID        uuid.UUID  `json:"sender_id"`
	Content         string     `json:"content"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
