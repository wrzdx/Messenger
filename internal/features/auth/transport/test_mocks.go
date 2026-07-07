package auth_transport_http

import (
	"context"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	"net/http"
)

type StubAuthService struct {
	CreateUserFn func(
		payload domain.RegisterUserPayload,
	) (
		domain.User,
		core_auth.AuthTokens,
		error,
	)

	LoginFn func(
		username string,
		password string,
	) (core_auth.AuthTokens, error)

	RefreshFn func(
		token string,
	) (core_auth.AuthTokens, error)
}

func (s *StubAuthService) Register(
	ctx context.Context,
	payload domain.RegisterUserPayload,
) (
	domain.User,
	core_auth.AuthTokens,
	error,
) {
	if s.CreateUserFn == nil {
		return domain.User{}, core_auth.AuthTokens{}, nil
	}
	return s.CreateUserFn(payload)
}

func (s *StubAuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (core_auth.AuthTokens, error) {
	if s.LoginFn == nil {
		return core_auth.AuthTokens{}, nil
	}
	return s.LoginFn(username, password)
}

func (s *StubAuthService) Refresh(
	ctx context.Context,
	token string,
) (core_auth.AuthTokens, error) {
	if s.LoginFn == nil {
		return core_auth.AuthTokens{}, nil
	}
	return s.RefreshFn(token)
}

type StubCookieManager struct {
	SetRefreshTokenFn func(
		w http.ResponseWriter,
		token string,
	)
	ClearRefreshTokenFn func(
		w http.ResponseWriter,
	)
}

func (m *StubCookieManager) SetRefreshToken(
	w http.ResponseWriter,
	token string,
) {
	if m.SetRefreshTokenFn != nil {
		m.SetRefreshTokenFn(w, token)
	}
}
func (m *StubCookieManager) ClearRefreshToken(
	w http.ResponseWriter,
) {
	if m.ClearRefreshTokenFn != nil {
		m.ClearRefreshTokenFn(w)
	}
}

func (m *StubCookieManager) GetRefreshToken(
	r *http.Request,
) (string, error) {
	return "", nil
}
