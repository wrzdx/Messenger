package auth_transport_http

import (
	"context"
	"messenger/internal/core/domain"
	"net/http"
)

type StubAuthService struct {
	CreateUserFn func(
		payload domain.RegisterUserPayload,
	) (
		domain.User,
		domain.TokenPair,
		error,
	)

	LoginFn func(
		username string,
		password string,
	) (domain.TokenPair, error)

	RefreshFn func(
		token string,
	) (domain.TokenPair, error)
}

func (s *StubAuthService) Register(
	ctx context.Context,
	payload domain.RegisterUserPayload,
) (
	domain.User,
	domain.TokenPair,
	error,
) {
	if s.CreateUserFn == nil {
		return domain.User{}, domain.TokenPair{}, nil
	}
	return s.CreateUserFn(payload)
}

func (s *StubAuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (domain.TokenPair, error) {
	if s.LoginFn == nil {
		return domain.TokenPair{}, nil
	}
	return s.LoginFn(username, password)
}

func (s *StubAuthService) Refresh(
	ctx context.Context,
	token string,
) (domain.TokenPair, error) {
	if s.LoginFn == nil {
		return domain.TokenPair{}, nil
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
