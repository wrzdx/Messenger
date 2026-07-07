package users_transport_http

import (
	core_auth "messenger/internal/core/auth"
	core_logger "messenger/internal/core/logger"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func (h *UsersHTTPHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	userID, ok := core_auth.UserIDFromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	if !ok {
		responseHandler.ErrorResponse(
			core_http_response.MapError(core_http_response.ErrMissingClaims),
		)
		return
	}
	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	responseHandler.NoContentResponse()
}
