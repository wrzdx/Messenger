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
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()

	if limit != nil && *limit < 0 {
		return nil, domain.ErrNegativeLimit
	}

	if offset != nil && *offset < 0 {
		return nil, domain.ErrNegativeOffset
	}

	query := `
	SELECT id, username, first_name, last_name, created_at, bio, password_hash
	FROM users
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: select users", err)
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
			&userModel.PasswordHash,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: scan users", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: next rows", err)
	}
	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
