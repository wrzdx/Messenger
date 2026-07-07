package users_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type UsersService struct {
	userRepository UsersRepository
}

type UsersRepository interface {
	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		userRepository: usersRepository,
	}
}
