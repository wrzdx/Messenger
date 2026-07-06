package users_postgres_repository

import (
	"context"
	"fmt"
	core_errors "messenger/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()
	query := `
	DELETE FROM users
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrorNotFound)
	}

	return nil
}
