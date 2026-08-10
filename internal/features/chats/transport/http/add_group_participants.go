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

func (h *ChatsHandler) AddGroupParticipants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	var request AddGroupParticipantsRequest
	if err := http_request.DecodeAndValidateRequestBody(r, &request); err != nil {
		sender.Error(err)
		return
	}
	ids := make([]uuid.UUID, 0, len(request.ParticipantIDs))
	for _, id := range request.ParticipantIDs {
		ids = append(ids, uuid.MustParse(id))
	}

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := uuid.Parse(chatIDStr)

	if err != nil {
		sender.Error(http_request.NewFieldError(map[string]string{
			"chat_id": "invalid uuid",
		}))
		return
	}
	result, err := h.chatsService.AddGroupParticipants(
		ctx,
		chats_service.AddGroupParticipantsCommand{
			GroupID:        chatID,
			RequesterID:    claims.UserID,
			ParticipantIDs: ids,
		},
	)
	if err != nil {
		sender.Error(err)
		return
	}

	response := make([]AddGroupParticipantItem, 0, len(result))
	for _, status := range result {
		response = append(response, AddGroupParticipantItem{
			UserID: status.UserID,
			Status: string(status.Status),
		})
	}
	sender.OK(http.StatusOK, response)
}

type AddGroupParticipantItem struct {
	UserID uuid.UUID `json:"user_id"`
	Status string    `json:"status"`
}

type AddGroupParticipantsResponse []AddGroupParticipantItem

type AddGroupParticipantsRequest struct {
	ParticipantIDs []string `json:"participant_ids" validate:"dive,uuid"`
}
