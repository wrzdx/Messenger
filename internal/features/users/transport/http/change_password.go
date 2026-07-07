package users_transport_http

import (
	"fmt"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=32" example:"password"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=32" example:"password"`
}

func (h *UsersHTTPHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request ChangePasswordRequest

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
}
