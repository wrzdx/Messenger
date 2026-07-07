package users_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

func (r *UsersRepository) ChangePassword(
	ctx context.Context,
	id uuid.UUID,
	passwordHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())
	defer cancel()

	query := `
	UPDATE users
	SET 
		password_hash=$1
	WHERE id=$2;`

	cmdTag, err := r.pool.Exec(ctx, query, passwordHash, id)
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
