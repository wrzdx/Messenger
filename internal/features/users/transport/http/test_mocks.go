package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"
)

type StubUsersService struct {
	CreateUserFn func(
		user domain.User,
		credentials domain.UserCredentials,
	) (domain.User, error)

	GetUsersFn func(
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUserFn func(
		id int,
	) (domain.User, error)

	DeleteUserFn func(
		id int,
	) error

	PatchUserFn func(
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func (s *StubUsersService) CreateUser(
	ctx context.Context,
	user domain.User,
	creds domain.UserCredentials,
) (domain.User, error) {
	return s.CreateUserFn(user, creds)
}

func (s *StubUsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	return s.GetUsersFn(limit, offset)
}

func (s *StubUsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	return s.GetUserFn(id)
}

func (s *StubUsersService) DeleteUser(
	id int,
) error {
	return s.DeleteUserFn(id)
}

func (s *StubUsersService) PatchUser(
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	return s.PatchUserFn(id, patch)
}
