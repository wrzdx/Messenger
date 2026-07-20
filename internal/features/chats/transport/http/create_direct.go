package chats_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"

	"github.com/google/uuid"
)

func (h *ChatsHandler) CreateDirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	var request CreateDirectRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}

	direct, isCreated, err := h.chatsService.CreateDirect(ctx, claims.UserID, request.PeerID)
	if err != nil {
		sender.Error(err)
		return
	}
	response := CreateDirectResponse(chatResponseFromDomain(direct.Chat))
	status := http.StatusOK
	if isCreated {
		status = http.StatusCreated
	}
	sender.OK(status, response)
}

type CreateDirectRequest struct {
	PeerID uuid.UUID `json:"peer_id" validate:"required"`
}

type CreateDirectResponse ChatResponse
