package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"
)

type StubUsersService struct {
	CreateUserFn func(
		ctx context.Context,
		user domain.User,
		credentials domain.UserCredentials,
	) (domain.User, error)

	GetUsersFn func(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUserFn func(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUserFn func(
		ctx context.Context,
		id int,
	) error

	PatchUserFn func(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func (s *StubUsersService) CreateUser(
	ctx context.Context,
	user domain.User,
	creds domain.UserCredentials,
) (domain.User, error) {
	return s.CreateUserFn(ctx, user, creds)
}

func (s *StubUsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	return s.GetUsersFn(ctx, limit, offset)
}

func (s *StubUsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	return s.GetUserFn(ctx, id)
}

func (s *StubUsersService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	return s.DeleteUserFn(ctx, id)
}

func (s *StubUsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	return s.PatchUserFn(ctx, id, patch)
}
