package auth_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/go-chi/chi/v5"
)

type AuthHTTPHandler struct {
	authService  AuthService
	cookieManger CookieManager
}

type AuthService interface {
	Register(
		ctx context.Context,
		payload domain.RegisterUserPayload,
	) (domain.User, domain.TokenPair, error)

	Login(
		ctx context.Context,
		username string,
		password string,
	) (domain.TokenPair, error)

	Refresh(
		ctx context.Context,
		token string,
	) (domain.TokenPair, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
	cookieManager CookieManager,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		cookieManger: cookieManager,
		authService:  authService,
	}
}

func (h *AuthHTTPHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/login", h.Login)
	router.Post("/register", h.Register)
	router.Post("/logout", h.Logout)
	router.Post("/refresh", h.Refresh)
	return router
}
