package auth_transport_http

import (
	"messenger/internal/core/domain"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	auth_service "messenger/internal/features/auth/service"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)
	var request RegisterRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}
	payload := auth_service.NewRegisterPayload(
		request.Username,
		request.FirstName,
		request.LastName,
		request.Bio,
		request.Password,
	)
	userDomain, tokens, err := h.authService.Register(ctx, payload)
	if err != nil {
		sender.Error(err)
		return
	}

	userResponse := UserResponse{
		userDomain.ID,
		userDomain.Username,
		userDomain.FirstName,
		userDomain.LastName,
		userDomain.CreatedAt,
		userDomain.Bio,
	}
	response := RegisterResponse{
		User:   userResponse,
		Access: tokens.Access,
	}

	h.cookieManger.SetRefreshToken(w, tokens.Refresh)
	sender.OK(http.StatusCreated, response)
}

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required" example:"qwerty"`
	FirstName string  `json:"first_name" validate:"required" example:"Ivan"`
	LastName  *string `json:"last_name"  example:"Ivanov"`
	Bio       *string `json:"bio"  example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
	Password  string  `json:"password" validate:"required" example:"password"`
}

func (r *RegisterRequest) Validate() map[string]string {
	fields := make(map[string]string)
	if err := domain.ValidateUsername(r.Username); err != nil {
		fields["username"] = err.Error()
	}

	if err := domain.ValidateFirstName(r.FirstName); err != nil {
		fields["first_name"] = err.Error()
	}

	if r.LastName != nil {
		if err := domain.ValidateLastName(*r.LastName); err != nil {
			fields["last_name"] = err.Error()
		}
	}

	if r.Bio != nil {
		if err := domain.ValidateBio(*r.Bio); err != nil {
			fields["bio"] = err.Error()
		}
	}

	if err := domain.ValidatePassword(r.Password); err != nil {
		fields["password"] = err.Error()
	}

	return fields
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"  example:"qwerty"`
	FirstName string    `json:"first_name"  example:"Ivan"`
	LastName  *string   `json:"last_name"  example:"Ivanov"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	Bio       *string   `json:"bio"  example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
}

type RegisterResponse struct {
	User   UserResponse `json:"user"`
	Access string       `json:"access_token"`
}
