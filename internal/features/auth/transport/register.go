package auth_transport_http

import (
	"fmt"
	"messenger/internal/core/domain"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required" example:"qwerty"`
	FirstName string  `json:"first_name" validate:"required" example:"Ivan"`
	LastName  *string `json:"last_name"  example:"Ivanov"`
	Bio       *string `json:"bio"  example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
	Password  string  `json:"password" validate:"required" example:"password"`
}

type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"  example:"qwerty"`
	FirstName string    `json:"first_name"  example:"Ivan"`
	LastName  *string   `json:"last_name"  example:"Ivanov"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	Bio       *string   `json:"bio"  example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request RegisterRequest
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
	payload := domain.NewRegisterUserPayload(
		request.Username,
		request.FirstName,
		request.LastName,
		request.Bio,
		request.Password,
	)
	userDomain, err := h.authService.Register(ctx, payload)
	if err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	response := RegisterResponse{
		userDomain.ID,
		userDomain.Username,
		userDomain.FirstName,
		userDomain.LastName,
		userDomain.CreatedAt,
		userDomain.Bio,
	}
	responseHandler.JSONResponse(response, http.StatusCreated)
}
