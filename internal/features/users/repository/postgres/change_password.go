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
	ctx, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()

	query := `
	UPDATE users
	SET password_hash=$1
	WHERE id=$2;`

	cmdTag, err := r.db.Exec(ctx, query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.NotFoundErr(
			domain.UserEntity,
			"id",
			id.String(),
		)
	}

	return nil
}
