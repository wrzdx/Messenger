package users_service

import (
	"context"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
)

type UsersService struct {
	hasher core_auth.PasswordHasher
	userRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
		passwordHash string,
	) (domain.User, error)

	// GetUsers(
	// 	ctx context.Context,
	// 	limit *int,
	// 	offset *int,
	// ) ([]domain.User, error)

	// GetUser(
	// 	ctx context.Context,
	// 	id int,
	// ) (domain.User, error)

	// DeleteUser(
	// 	ctx context.Context,
	// 	id int,
	// ) error

	// PatchUser(
	// 	ctx context.Context,
	// 	id int,
	// 	user domain.User,
	// ) (domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository,
	hasher core_auth.PasswordHasher,
) *UsersService {
	return &UsersService{
		userRepository: usersRepository,
		hasher: hasher,
	}
}
