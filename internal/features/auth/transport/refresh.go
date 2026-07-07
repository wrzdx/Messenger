package auth_transport_http

import (
	"errors"
	core_logger "messenger/internal/core/logger"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type RefreshResponse struct {
	Access string `json:"access"`
}

func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	token, err := h.cookieManger.GetRefreshToken(r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			responseHandler.ErrorResponse(
				core_http_response.Error{
					Error:   err,
					Status:  http.StatusUnauthorized,
					Message: "Missing credentials",
				},
			)
			return
		}
		responseHandler.ErrorResponse(
			core_http_response.Error{
				Error:   err,
				Status:  http.StatusUnauthorized,
				Message: "Failed to get refresh token",
			},
		)
	}

	tokens, err := h.authService.Refresh(ctx, token)
	if err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}
	response := RefreshResponse{
		Access: tokens.Access,
	}

	h.cookieManger.SetRefreshToken(w, tokens.Refresh)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
