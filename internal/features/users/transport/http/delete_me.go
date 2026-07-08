package users_transport_http

import (
	auth "messenger/internal/core/auth"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func (h *UsersHTTPHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	userID := auth.MustUserIDFromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(http_response.MapError(err))
		return
	}

	responseHandler.NoContentResponse()
}
