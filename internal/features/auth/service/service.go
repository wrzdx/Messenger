package auth_service

import (
	"context"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
)

type AuthService struct {
	userRepository UsersRepository
	hasher         core_auth.PasswordHasher
	jwtProvider    core_auth.JWTProvider
}

type UsersRepository interface {
	GetUserAuth(
		ctx context.Context,
		username string,
	) (domain.UserAuth, error)

	// UpdateUserPasswordHash(
	// 	ctx context.Context,
	// 	username string,
	// 	passwordHash string,
	// ) error
}

func NewUsersService(
	usersRepository UsersRepository,
	hasher core_auth.PasswordHasher,
	jwtProvider core_auth.JWTProvider,
) *AuthService {
	return &AuthService{
		userRepository: usersRepository,
		hasher:         hasher,
		jwtProvider:    jwtProvider,
	}
}
