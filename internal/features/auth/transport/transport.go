package auth_transport_http

import (
	"context"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	auth_service "messenger/internal/features/auth/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService  AuthService
	cookieManger CookieManager
}

type AuthService interface {
	Register(
		ctx context.Context,
		payload auth_service.RegisterPayload,
	) (domain.User, auth.TokenPair, error)

	Login(
		ctx context.Context,
		username string,
		password string,
	) (auth.TokenPair, error)

	Refresh(
		ctx context.Context,
		token string,
	) (auth.TokenPair, error)
}

type CookieManager interface {
	SetRefreshToken(
		w http.ResponseWriter,
		token string,
	)

	ClearRefreshToken(
		w http.ResponseWriter,
	)

	GetRefreshToken(
		r *http.Request,
	) (string, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
	cookieManager CookieManager,
) *AuthHandler {
	return &AuthHandler{
		cookieManger: cookieManager,
		authService:  authService,
	}
}

func (h *AuthHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/login", h.Login)
	router.Post("/register", h.Register)
	router.Post("/logout", h.Logout)
	router.Post("/refresh", h.Refresh)
	return router
}
