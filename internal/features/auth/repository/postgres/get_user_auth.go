package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors.go"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetUserAuth(
	ctx context.Context,
	username string,
) (domain.UserAuth, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()

	query := `
	SELECT id, password_hash FROM users 
	WHERE username=$1;
	`
	var userAuthModel UserAuthModel
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&userAuthModel.UserID,
		&userAuthModel.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserAuth{}, fmt.Errorf(
				"user with username='%s': %w",
				username,
				core_errors.ErrorNotFound,
			)
		}
		return domain.UserAuth{}, fmt.Errorf("scan error: %w", err)
	}

	userAuthDomain := UserAuthDomainFromModel(userAuthModel)
	return userAuthDomain, nil
}
