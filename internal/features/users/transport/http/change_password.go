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
	userID := core_auth.MustUserIDFromContext(ctx)

	var request ChangePasswordRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		err = fmt.Errorf(
			"%v: %w",
			err,
			core_http_response.ErrInvalidArgument,
		)
		responseHandler.ErrorResponse(
			core_http_response.Error{
				Error:   err,
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	if err := h.usersService.ChangePassword(
		ctx,
		userID,
		request.OldPassword,
		request.NewPassword,
	); err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	responseHandler.NoContentResponse()
}
