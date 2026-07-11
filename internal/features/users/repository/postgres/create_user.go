package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	postgres "messenger/internal/core/repository/postgres"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()

	query := `
	INSERT INTO users (id, username, first_name, last_name, created_at, deleted_at, bio, password_hash)
	VALUES ($1, $2,$3,$4,$5,$6, $7,$8) 
	RETURNING 
		id,
		username,
		first_name,
		last_name,
		created_at,
		deleted_at,
		bio,
		password_hash;
	`
	var userModel UserModel
	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.DeletedAt,
		user.Bio,
		user.PasswordHash,
	).Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.CreatedAt,
		&user.DeletedAt,
		&userModel.Bio,
		&userModel.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrViolatesUnique) {
			return domain.User{}, domain.ErrAlreadyExists
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := UserDomainFromModel(userModel)
	return userDomain, nil
}
