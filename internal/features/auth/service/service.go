package auth_service

import (
	"context"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	usersRepository UsersRepository
	hasher          Hasher
	jwtProvider     JWTProvider
}

type UsersRepository interface {
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	GetUserByUsername(
		ctx context.Context,
		username string,
	) (domain.User, error)

	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type JWTProvider interface {
	GenerateTokens(id uuid.UUID) (core_auth.AuthTokens, error)
	ParseToken(token string) (core_auth.Claims, error)
}

func NewAuthService(
	userRepository UsersRepository,
	hasher Hasher,
	jwtProvider JWTProvider,
) *AuthService {
	return &AuthService{
		usersRepository: userRepository,
		hasher:          hasher,
		jwtProvider:     jwtProvider,
	}
}
