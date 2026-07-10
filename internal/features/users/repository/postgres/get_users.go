package users_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil {
		if err := domain.ValidateLimit(*limit); err != nil {
			return nil, err
		}
	}

	if offset != nil {
		if err := domain.ValidateOffset(*offset); err != nil {
			return nil, err
		}
	}
	ctx, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()

	query := `
	SELECT *
	FROM users
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
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
	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
