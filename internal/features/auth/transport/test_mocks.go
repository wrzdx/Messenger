package auth_transport_http

import (
	"context"
	"messenger/internal/core/domain"
)

type StubAuthService struct {
	CreateUserFn func(
		payload domain.RegisterUserPayload,
	) (domain.User, error)

	LoginFn func(
		username string,
		password string,
	) (domain.Token, domain.Token, error)
}

func (s *StubAuthService) Register(
	ctx context.Context,
	payload domain.RegisterUserPayload,
) (domain.User, error) {
	return s.CreateUserFn(payload)
}

func (s *StubAuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (domain.Token, domain.Token, error) {
	return s.LoginFn(username, password)
}
