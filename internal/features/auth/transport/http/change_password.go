package auth_transport_http

import (
	"errors"
	"messenger/internal/core/auth"
	core_context "messenger/internal/core/context"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"password"`
	NewPassword     string `json:"new_password" validate:"required" example:"password"`
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	var request ChangePasswordRequest

	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}

	if err := h.authService.ChangePassword(
		ctx,
		claims.UserID,
		request.CurrentPassword,
		request.NewPassword,
	); err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			h.cookieManager.ClearRefreshToken(w)
		}
		sender.Error(err)
		return
	}
	h.cookieManager.ClearRefreshToken(w)
	sender.OK(http.StatusNoContent, nil)
}
