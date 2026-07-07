package auth_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepository AuthRepository
	hasher         Hasher
	jwtProvider    JWTProvider
}

type AuthRepository interface {
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
	GenerateAccessToken(id uuid.UUID) (domain.Token, error)
	GenerateRefreshToken(id uuid.UUID) (domain.Token, error)
	ParseToken(token string) (domain.Claims, error)
}

func NewAuthService(
	authRepository AuthRepository,
	hasher Hasher,
	jwtProvider JWTProvider,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		hasher:         hasher,
		jwtProvider:    jwtProvider,
	}
}
