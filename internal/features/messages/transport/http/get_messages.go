package messages_transport_http

import (
	"fmt"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	messages_service "messenger/internal/features/messages/service"
	"net/http"
	"strconv"

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
		sender.Error(fmt.Errorf("invalid chat id: %w", http_request.ErrInvalidRequest))
		return
	}

	queryParams := r.URL.Query()
	cursorStr := queryParams.Get("cursor")
	cursor, err := decodeMessageCursor(cursorStr)
	if err != nil {
		sender.Error(err)
		return
	}

	var limit int
	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			sender.Error(fmt.Errorf(
				"invalid limit query param: %w",
				http_request.ErrInvalidRequest,
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
	nextCursor, err := encodeMessageCursor(page.NextCursor)
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
