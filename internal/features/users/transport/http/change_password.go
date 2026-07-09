package users_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"password"`
	NewPassword string `json:"new_password" validate:"required" example:"password"`
}

func (r *ChangePasswordRequest) Validate() map[string]string {
	fields := make(map[string]string)

	if err := domain.ValidatePassword(r.NewPassword); err != nil {
		fields["new_password"] = err.Error()
	}

	return fields
}

func (h *UsersHTTPHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)
	claims := core_context.ClaimsRequired(ctx)

	var request ChangePasswordRequest

	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}

	if err := h.usersService.ChangePassword(
		ctx,
		claims.UserID,
		request.OldPassword,
		request.NewPassword,
	); err != nil {
		sender.Error(err)
		return
	}

	sender.OK(http.StatusNoContent, nil)
}
