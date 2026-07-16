package users_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

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

	UpdateUserProfile(
		ctx context.Context,
		id uuid.UUID,
		profile domain.UserProfile,
	) error
}
