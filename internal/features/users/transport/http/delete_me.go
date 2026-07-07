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
	claims, ok := core_auth.ClaimsFromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	if !ok {
		responseHandler.ErrorResponse(
			http.StatusUnauthorized,
			core_http_response.ErrMissingClaims,
		)
		return
	}
	if err := h.usersService.DeleteUser(ctx, claims.UserID); err != nil {
		statusCode := core_http_response.MapDomainErrorToStatusCode(err)
		responseHandler.ErrorResponse(statusCode, err)
		return
	}

	responseHandler.NoContentResponse()
}
