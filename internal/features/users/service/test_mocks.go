package users_service

import (
	"context"
	"messenger/internal/core/domain"
)

type StubUsersRepository struct {
	CreateUserFn func(
		ctx context.Context,
		user domain.User,
		passwordHash string,
	) (domain.User, error)
}

func (s *StubUsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
	passwordHash string,
) (domain.User, error) {
	return s.CreateUserFn(ctx, user, passwordHash)
}

type StubHasher struct {
	HashFn    func(password string) ([]byte, error)
	CompareFn func(hash, password string) error
}

func (h *StubHasher) Hash(password string) ([]byte, error) {
	return h.HashFn(password)
}
func (h *StubHasher) Compare(hash, password string) error {
	return h.CompareFn(hash, password)
}
