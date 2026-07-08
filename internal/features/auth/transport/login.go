package auth_transport_http

import (
	"fmt"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"qwerty"`
	Password string `json:"password" validate:"required" example:"password"`
}

type LoginResponse struct {
	Access string `json:"access"`
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)
	var request LoginRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			http_response.Error{
				Error: fmt.Errorf(
					"%v: %w",
					err,
					http_response.ErrInvalidArgument,
				),
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	tokens, err := h.authService.Login(
		ctx,
		request.Username,
		request.Password,
	)
	if err != nil {
		responseHandler.ErrorResponse(http_response.MapError(err))
		return
	}
	response := LoginResponse{
		Access: tokens.Access,
	}

	h.cookieManger.SetRefreshToken(w, tokens.Refresh)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
