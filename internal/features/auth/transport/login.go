package auth_transport_http

import (
	"messenger/internal/core/domain"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" example:"qwerty"`
	Password string `json:"password" validate:"required,min=8,max=32" example:"password"`
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
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userCredentials := domain.NewCredentials(request.Username, request.Password)
	refresh, access, err := h.authService.Login(ctx, userCredentials)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}
	response := LoginResponse{
		Access: access.Token,
	}
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Token,
		Secure:   h.secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refresh.Expires,
	}
	http.SetCookie(w, cookie)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
