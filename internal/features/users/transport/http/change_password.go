package users_transport_http

import (
	"fmt"
	auth "messenger/internal/core/auth"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"password"`
	NewPassword string `json:"new_password" validate:"required" example:"password"`
}

func (h *UsersHTTPHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)
	userID := auth.MustUserIDFromContext(ctx)

	var request ChangePasswordRequest

	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		err = fmt.Errorf(
			"%v: %w",
			err,
			http_response.ErrInvalidArgument,
		)
		responseHandler.ErrorResponse(
			http_response.Error{
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
		responseHandler.ErrorResponse(http_response.MapError(err))
		return
	}

	responseHandler.NoContentResponse()
}
