package messages_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_cursor "messenger/internal/core/transport/http/cursor"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	messages_service "messenger/internal/features/messages/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *MessagesHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := uuid.Parse(chatIDStr)

	if err != nil {
		sender.Error(http_request.NewFieldError(map[string]string{
			"chat_id": "invalid uuid",
		}))
		return
	}

	queryParams := r.URL.Query()
	var cursor *messages_service.MessageCursor

	cursorPayload, err := http_cursor.DecodeAndValidate[messageCursorPayload](
		queryParams.Get("cursor"),
	)

	if err != nil {
		sender.Error(err)
		return
	}
	if cursorPayload != nil {
		cursor = &messages_service.MessageCursor{
			MessageID: uuid.MustParse(cursorPayload.MessageID),
			CreatedAt: cursorPayload.CreatedAt,
		}
	}

	var limit int
	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			sender.Error(http_request.NewFieldError(
				map[string]string{"limit": "invalid limit query param"},
			))
			return
		}
	}

	page, err := h.messagesService.GetMessages(
		ctx,
		claims.UserID,
		messages_service.GetMessagesQuery{
			ChatID: chatID,
			Before: cursor,
			Limit:  limit,
		},
	)
	if err != nil {
		sender.Error(err)
		return
	}
	var responseCursorPayload *messageCursorPayload
	if page.NextCursor != nil {
		responseCursorPayload = &messageCursorPayload{
			MessageID: page.NextCursor.MessageID.String(),
			CreatedAt: page.NextCursor.CreatedAt,
		}
	}

	nextCursor, err := http_cursor.Encode(responseCursorPayload)
	if err != nil {
		sender.Error(err)
		return
	}
	sender.OK(http.StatusOK, GetMessagesResponse{
		Messages:   messageResponsesFromDomain(page.Messages),
		NextCursor: nextCursor,
	})
}

type GetMessagesResponse struct {
	Messages   []MessageResponse `json:"messages"`
	NextCursor *string           `json:"next_cursor"`
}

type messageCursorPayload struct {
	MessageID string    `json:"message_id" validate:"required,uuid"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}
