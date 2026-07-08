package users_transport_http

import (
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"

	"github.com/google/uuid"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)

	userID, err := http_request.GetPathValue[uuid.UUID](r, "id")
	if err != nil {
		sender.Error(err)
		return
	}
	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		sender.Error(err)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	sender.OK(http.StatusOK, response)
}
