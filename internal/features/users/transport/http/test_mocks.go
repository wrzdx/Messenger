package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"
)
type StubUsersService struct {
    Called      bool
    GotUser     domain.User
    GotCreds    domain.UserCredentials

    ReturnUser  domain.User
    ReturnError error
}

func (s *StubUsersService) CreateUser(
    ctx context.Context,
    user domain.User,
    creds domain.UserCredentials,
) (domain.User, error) {
    s.Called = true
    s.GotUser = user
    s.GotCreds = creds

    return s.ReturnUser, s.ReturnError
}
