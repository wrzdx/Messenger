package users_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type UsersService struct {
	userRepository UsersRepository
	hasher         Hasher
}

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type UsersRepository interface {
	GetUsers(
		ctx context.Context,
		pagination domain.Pagination,
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

	ChangePassword(
		ctx context.Context,
		id uuid.UUID,
		newPasswordHash string,
	) error
}

func NewUsersService(
	usersRepository UsersRepository,
	hasher Hasher,
) *UsersService {
	return &UsersService{
		userRepository: usersRepository,
		hasher:         hasher,
	}
}
