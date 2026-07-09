package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	postgres "messenger/internal/core/repository/postgres"
)

func (r *UsersRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()
	query := `
	SELECT id, username, first_name, last_name, created_at, bio, password_hash
	FROM users
	WHERE username=$1;
	`

	row := r.db.QueryRow(ctx, query, username)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.CreatedAt,
		&userModel.Bio,
		&userModel.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.User{}, domain.NotFoundErr(
				domain.UserEntity,
				"username",
				username,
			)
		}

		return domain.User{}, fmt.Errorf(
			"scan user by username %q: %w",
			username,
			err,
		)
	}

	userDomain := UserDomainFromModel(userModel)

	return userDomain, nil
}
