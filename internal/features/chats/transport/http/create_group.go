package chats_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	chats_service "messenger/internal/features/chats/service"
	"net/http"

	"github.com/google/uuid"
)

func (h *ChatsHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	var request CreateGroupRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}
	ids := make([]uuid.UUID, 0, len(request.ParticipantIDs))
	for _, id := range request.ParticipantIDs {
		ids = append(ids, uuid.MustParse(id))
	}
	group, err := h.chatsService.CreateGroup(
		ctx,
		claims.UserID,
		chats_service.CreateGroupCommand{
			Title:          request.Title,
			ParticipantIDs: ids,
		},
	)
	if err != nil {
		sender.Error(err)
		return
	}

	response := CreateGroupResponse(chatResponseFromDomain(group.Chat))
	sender.OK(http.StatusCreated, response)
}

type CreateGroupResponse ChatResponse

type CreateGroupRequest struct {
	Title          string   `json:"title" validate:"required"`
	ParticipantIDs []string `json:"participant_ids" validate:"dive,uuid"`
}
