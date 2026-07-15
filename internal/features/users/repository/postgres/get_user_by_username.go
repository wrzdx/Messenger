package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"

	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	db := postgres.GetExecutor(ctx, r.db)
	query := `
	SELECT *
	FROM users
	WHERE lower(username)=lower($1);
	`

	row := db.QueryRow(ctx, query, username)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.CreatedAt,
		&userModel.DeletedAt,
		&userModel.Bio,
		&userModel.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf(
			"scan user by username %q: %w",
			username,
			err,
		)
	}

	userDomain, err := UserDomainFromModel(userModel)
	if err != nil {
		return domain.User{}, err
	}
	return userDomain, nil
}
