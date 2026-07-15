package users_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	pagination domain.Pagination,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	db := postgres.GetExecutor(ctx, r.db)

	query := `
	SELECT *
	FROM users
	LIMIT $1
	OFFSET $2;
	`

	rows, err := db.Query(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()
	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}
	userDomains, err := userDomainsFromModels(userModels)
	if err != nil {
		return nil, err
	}
	return userDomains, nil
}
