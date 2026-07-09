package users_transport_http

import (
	core_context "messenger/internal/core/context"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func (h *UsersHTTPHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	claims := core_context.ClaimsRequired(ctx)
	sender := http_response.NewHTTPSender(log, w)

	user, err := h.usersService.GetUser(ctx, claims.UserID)
	if err != nil {
		sender.Error(err)
	}

	response := GetUserResponse(userDTOFromDomain(user))
	sender.OK(http.StatusOK, response)
}
