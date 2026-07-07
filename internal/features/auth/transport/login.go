package auth_transport_http

import (
	"fmt"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
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
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
	var request LoginRequest
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

	tokens, err := h.authService.Login(
		ctx,
		request.Username,
		request.Password,
	)
	if err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}
	response := LoginResponse{
		Access: tokens.Access,
	}

	h.cookieManger.SetRefreshToken(w, tokens.Refresh)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
