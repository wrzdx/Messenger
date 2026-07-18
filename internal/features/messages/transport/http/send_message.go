package messages_transport_http

import (
	"net/http"

	"github.com/google/uuid"
)

func (h *MessagesHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// log := logger.FromContext(ctx)
	// sender := http_response.NewHTTPSender(log, w, errorMapper)
	// claims := core_context.ClaimsRequired(ctx)

	// var request SendMessageRequest
	// if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
	// 	sender.Error(err)
	// }

}

type SendMessageRequest struct {
	ClientMessageID uuid.UUID `json:"client_message_id" validate:"required"`
	Content         string    `json:"content" validate:"required"`
}
