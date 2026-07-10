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
	UPDATE users
	SET 
		username='deleted_' || substr(md5(id::text), 1, 16),
		first_name='Deleted Account',
		last_name=NULL,
		deleted_at=NOW(),
		bio=NULL,
		password_hash=''
	WHERE id=$1 AND deleted_at IS NULL;
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
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
