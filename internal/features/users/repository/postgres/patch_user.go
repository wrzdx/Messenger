package users_postgres_repository

import (
	"context"
	"messenger/internal/core/domain"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	return domain.User{}, nil
}
