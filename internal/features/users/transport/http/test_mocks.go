package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type StubUsersService struct {
	CreateUserFn func(
		payload domain.RegisterUserPayload,
	) (domain.User, error)

	GetUsersFn func(
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUserFn func(
		id uuid.UUID,
	) (domain.User, error)

	DeleteUserFn func(
		id uuid.UUID,
	) error

	PatchUserFn func(
		id uuid.UUID,
		patch domain.UserPatch,
	) (domain.User, error)

	ChangePasswordFn func(
		id uuid.UUID,
		old_password string,
		new_password string,
	) error
}

func (s *StubUsersService) CreateUser(
	ctx context.Context,
	payload domain.RegisterUserPayload,
) (domain.User, error) {
	return s.CreateUserFn(payload)
}

func (s *StubUsersService) GetUsers(
	ctx context.Context,
	pagination domain.Pagination,
) ([]domain.User, error) {
	return s.GetUsersFn(pagination.Limit, pagination.Offset)
}

func (s *StubUsersService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	return s.GetUserFn(id)
}

func (s *StubUsersService) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.DeleteUserFn(id)
}

func (s *StubUsersService) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	patch domain.UserPatch,
) (domain.User, error) {
	return s.PatchUserFn(id, patch)
}

func (s *StubUsersService) ChangePassword(
	ctx context.Context,
	id uuid.UUID,
	old_password string,
	new_password string,
) error {
	return s.ChangePasswordFn(id, old_password, new_password)
}
