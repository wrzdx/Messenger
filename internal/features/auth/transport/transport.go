package auth_transport_http

import (
	"context"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	"net/http"

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
	) (domain.User, core_auth.AuthTokens, error)

	Login(
		ctx context.Context,
		username string,
		password string,
	) (core_auth.AuthTokens, error)

	Refresh(
		ctx context.Context,
		token string,
	) (core_auth.AuthTokens, error)
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
