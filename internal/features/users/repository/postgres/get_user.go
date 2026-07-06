package users_postgres_repository

import (
	"context"
	"messenger/internal/core/domain"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	return domain.User{}, nil
}
