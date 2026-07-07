package users_transport_http

import (
	core_auth "messenger/internal/core/auth"
	core_logger "messenger/internal/core/logger"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func (h *UsersHTTPHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	claims, ok := core_auth.ClaimsFromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	if !ok {
		responseHandler.ErrorResponse(
			core_http_response.MapError(core_http_response.ErrMissingClaims),
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(
			core_http_response.MapError(err),
		)
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
