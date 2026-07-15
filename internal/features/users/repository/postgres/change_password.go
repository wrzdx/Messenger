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

func (r *UsersRepository) ChangePassword(
	ctx context.Context,
	id uuid.UUID,
	passwordHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	db := postgres.GetExecutor(ctx, r.db)

	query := `
	UPDATE users
	SET password_hash=$1
	WHERE id=$2
	RETURNING id;`

	err := db.QueryRow(ctx, query, passwordHash, id).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
