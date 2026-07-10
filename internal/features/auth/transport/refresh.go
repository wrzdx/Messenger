package auth_transport_http

import (
	"fmt"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type RefreshResponse struct {
	Access string `json:"access"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)

	token, err := h.cookieManger.GetRefreshToken(r)
	if err != nil {
		sender.Error(fmt.Errorf("get refresh token: %w", err))
		return
	}

	tokens, err := h.authService.Refresh(ctx, token)
	if err != nil {
		sender.Error(err)
		return
	}
	response := RefreshResponse{
		Access: tokens.Access,
	}

	h.cookieManger.SetRefreshToken(w, tokens.Refresh)
	sender.OK(http.StatusCreated, response)
}
