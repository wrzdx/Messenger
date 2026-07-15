package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	db := postgres.GetExecutor(ctx, r.db)
	query := `
	UPDATE users
	SET 
		username='deleted_' || substr(md5(id::text), 1, 16),
		first_name='Deleted Account',
		last_name=NULL,
		deleted_at=NOW(),
		bio=NULL
	WHERE id=$1 AND deleted_at IS NULL
	RETURNING id;
	`

	err := db.QueryRow(ctx, query, id).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("delete user query: %w", err)
	}

	return nil
}
