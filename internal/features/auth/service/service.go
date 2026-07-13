package auth_service

import (
	"context"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	usersRepository UsersRepository
	hasher          Hasher
	tokenProvider   TokenProvider
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

type 

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenProvider interface {
	GenerateAccessToken(
		accessTokenClaims auth.AccessTokenClaims,
		tokenLifetime auth.TokenLifetime,
	) (string, error)
	GenerateRefreshToken(
		refreshTokenClaims auth.RefreshTokenClaims,
		tokenLifetime auth.TokenLifetime,
	) (string, error)
	ParseAccessToken(tokenStr string) (auth.AccessTokenClaims, error)
	ParseRefreshToken(tokenStr string) (auth.RefreshTokenClaims, error)
}

func NewAuthService(
	userRepository UsersRepository,
	hasher Hasher,
	tokenProvider TokenProvider,
) *AuthService {
	return &AuthService{
		usersRepository: userRepository,
		hasher:          hasher,
		tokenProvider:   tokenProvider,
	}
}
