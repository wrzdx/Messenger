package auth_service

import (
	"context"
	auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	usersRepository UsersRepository
	hasher          auth.Hasher
	tokenService    auth.TokenService
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

func NewAuthService(
	userRepository UsersRepository,
	hasher auth.Hasher,
	tokenService auth.TokenService,
) *AuthService {
	return &AuthService{
		usersRepository: userRepository,
		hasher:          hasher,
		tokenService:    tokenService,
	}
}
