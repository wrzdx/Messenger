package users_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()
	query := `
	DELETE FROM users
	WHERE id=$1;
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w",
			id,
			domain.ErrUserNotFound,
		)
	}

	return nil
}
