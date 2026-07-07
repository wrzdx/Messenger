package auth_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type StubAuthRepository struct {
	CreateUserFn func(
		user domain.User,
	) (domain.User, error)

	GetUserFn func(
		username string,
	) (domain.User, error)
}

func (s *StubAuthRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return s.CreateUserFn(user)
}

func (s *StubAuthRepository) GetUser(
	ctx context.Context,
	username string,
) (domain.User, error) {
	return s.GetUserFn(username)
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

type StubJWTProvider struct {
}

func (p *StubJWTProvider) GenerateAccessToken(id uuid.UUID) (domain.Token, error) {
	return domain.Token{}, nil
}
func (p *StubJWTProvider) GenerateRefreshToken(id uuid.UUID) (domain.Token, error) {
	return domain.Token{}, nil
}
func (p *StubJWTProvider) ParseToken(token string) (domain.Claims, error) {
	return domain.Claims{}, nil
}
