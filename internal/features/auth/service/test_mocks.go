package auth_service

import (
	"context"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type StubAuthRepository struct {
	CreateUserFn func(
		user domain.User,
	) (domain.User, error)

	GetUserByUsernameFn func(
		username string,
	) (domain.User, error)
}

func (s *StubAuthRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return s.CreateUserFn(user)
}

func (s *StubAuthRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (domain.User, error) {
	return s.GetUserByUsernameFn(username)
}

type StubHasher struct {
	HashFn    func(password string) (string, error)
	CompareFn func(hash, password string) error
}

func (h *StubHasher) Hash(password string) (string, error) {
	return h.HashFn(password)
}
func (h *StubHasher) Compare(hash, password string) error {
	return h.CompareFn(hash, password)
}

type StubJWTProvider struct {
	GenerateTokensFn func(id uuid.UUID) (core_auth.AuthTokens, error)
	ParseTokenFn     func(token string) (core_auth.Claims, error)
}

func (p *StubJWTProvider) GenerateTokens(id uuid.UUID) (core_auth.AuthTokens, error) {
	if p.GenerateTokensFn != nil {
		return p.GenerateTokensFn(id)
	}
	return core_auth.AuthTokens{}, nil
}

func (p *StubJWTProvider) ParseToken(token string) (core_auth.Claims, error) {
	if p.ParseTokenFn != nil {
		return p.ParseTokenFn(token)
	}
	return core_auth.Claims{}, nil
}
