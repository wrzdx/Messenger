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
	Login(
		ctx context.Context,
		credentials domain.UserCredentials,
	) (domain.Token, domain.Token, error)
}

func NewUsersHTTPHandler(userService AuthService, secure bool) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		secure:      secure,
		authService: userService,
	}
}

func (h *AuthHTTPHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/login", h.Login)
	return router
}
