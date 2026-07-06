package users_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	query := `
	SELECT id, username, first_name, last_name, created_at, bio 
	FROM users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
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
			&userModel.Bio,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}
	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
