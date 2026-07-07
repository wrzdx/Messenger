package auth_transport_http

import (
	"fmt"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" example:"qwerty"`
	Password string `json:"password" example:"password"`
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
			http.StatusBadRequest,
			fmt.Errorf("%w: %v", core_http_response.ErrInvalidArgument, err),
		)
		return
	}

	refresh, access, err := h.authService.Login(
		ctx,
		request.Username,
		request.Password,
	)
	if err != nil {
		statusCode := core_http_response.MapDomainErrorToStatusCode(err)
		responseHandler.ErrorResponse(statusCode, err)
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
