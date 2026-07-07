package users_transport_http

import (
	"fmt"
	core_auth "messenger/internal/core/auth"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"password"`
	NewPassword string `json:"new_password" validate:"required" example:"password"`
}

func (h *UsersHTTPHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.ClaimsFromContext(ctx)
	var request ChangePasswordRequest

	if !ok {
		responseHandler.ErrorResponse(
			core_http_response.MapError(core_http_response.ErrMissingClaims),
		)
		return
	}

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			core_http_response.MapError(
				fmt.Errorf(
					"%v: %w",
					err,
					core_http_response.ErrInvalidArgument,
				),
			),
		)
		return
	}

	if err := h.usersService.ChangePassword(
		ctx,
		claims.UserID,
		request.OldPassword,
		request.NewPassword,
	); err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	responseHandler.NoContentResponse()
}
