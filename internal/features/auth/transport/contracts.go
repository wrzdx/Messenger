package auth_transport_http

import (
	"context"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	auth_service "messenger/internal/features/auth/service"
	"net/http"
)

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

	Logout(ctx context.Context, refreshToken string) error
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
