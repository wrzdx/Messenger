package chats_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	chats_service "messenger/internal/features/chats/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *ChatsHandler) RemoveGroupParticipant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	var request RemoveGroupParticipantRequest
	if err := http_request.DecodeAndValidateRequestBody(r, &request); err != nil {
		sender.Error(err)
		return
	}

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		sender.Error(http_request.NewFieldError(
			map[string]string{
				"chat_id": "invalid uuid",
			},
		))
		return
	}
	targetID := uuid.MustParse(request.TargetID)

	if err := h.chatsService.RemoveGroupParticipant(
		ctx,
		chats_service.RemoveGroupParticipantCommand{
			GroupID:     chatID,
			RequesterID: claims.UserID,
			TargetID:    targetID,
		}); err != nil {
		sender.Error(err)
		return
	}

	sender.OK(http.StatusNoContent, nil)
}

type RemoveGroupParticipantRequest struct {
	TargetID string `json:"target_id" validate:"uuid"`
}
