package domain

import (
	"context"

	"github.com/google/uuid"
)

type UsersService interface {
	CreateUser(
		ctx context.Context,
		payload RegisterUserPayload,
	) (User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		patch UserPatch,
	) (User, error)
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user User,
	) (User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		user User,
	) (User, error)
}
