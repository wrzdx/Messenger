package chats_transport_http

import (
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"

	"github.com/google/uuid"
)

func (h *ChatsHandler) CreateDirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)

	var request CreateDirectRequest
	if err:= http_request.DecodeAndValidateRequest(r, request); err!= nil {
		sender.Error(err)
		return 
	}
}


type CreateDirectRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}