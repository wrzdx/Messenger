package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
	passwordHash string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()

	query := `
	INSERT INTO users (username, first_name, last_name, created_at, bio, password_hash)
	VALUES ($1, $2,$3,$4,$5,$6) 
	RETURNING id, username, first_name, last_name, created_at, bio;
	`
	var userModel UserModel
	err := r.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.Bio,
		passwordHash,
	).Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.CreatedAt,
		&userModel.Bio,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.User{}, fmt.Errorf(
				"%v: user with username=%s already exists: %w",
				err,
				user.Username,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userAuthDomain := UserDomainFromModel(userModel)
	return userAuthDomain, nil
}
