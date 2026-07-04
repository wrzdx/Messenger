package users_service

import (
	"context"
	"messenger/internal/core/domain"
)

type StubUsersRepository struct {
	Called     bool
	GotUser    domain.User
	GotPswHash string

	ReturnUser  domain.User
	ReturnError error
}

func (s *StubUsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
	passwordHash string,
) (domain.User, error) {
	s.Called = true
	s.GotUser = user
	s.GotPswHash = passwordHash

	return s.ReturnUser, s.ReturnError
}

type StubHasher struct {
	Called      bool
	GotPassword string

	ReturnHash  []byte
	ReturnError error
}

func (h *StubHasher) Hash(password string) ([]byte, error) {
	h.Called = true
	h.GotPassword = password
	return h.ReturnHash, h.ReturnError
}
func (h *StubHasher) Compare(hash, password string) error {
	return nil
}
