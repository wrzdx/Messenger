package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()

	query := `
	INSERT INTO users (id, username, first_name, last_name, created_at, bio, password_hash)
	VALUES ($1, $2,$3,$4,$5,$6, $7) 
	RETURNING id, username, first_name, last_name, created_at, bio, password_hash;
	`
	var userModel UserModel
	err := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.Bio,
		user.PasswordHash,
	).Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.CreatedAt,
		&userModel.Bio,
		&userModel.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.User{}, fmt.Errorf(
				"%w: user with username=%s already exists",
				err,
				user.Username,
			)
		}
		return domain.User{}, fmt.Errorf("%w: scan error", err)
	}

	userDomain := UserDomainFromModel(userModel)
	return userDomain, nil
}
