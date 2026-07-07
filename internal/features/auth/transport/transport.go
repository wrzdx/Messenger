package auth_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/go-chi/chi/v5"
)

type AuthHTTPHandler struct {
	authService AuthService
	secure      bool
}

type AuthService interface {
	Register(
		ctx context.Context,
		payload domain.RegisterUserPayload,
	) (domain.User, error)
	Login(
		ctx context.Context,
		username string,
		password string,
	) (domain.Token, domain.Token, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
	secure bool,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		secure:      secure,
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/login", h.Login)
	router.Post("/register", h.Register)
	return router
}
