package core_test_utils

import (
	"context"
	core_logger "messenger/internal/core/logger"
	core_postgres "messenger/internal/core/repository/postgres"
	"testing"
)

func ResetDB(t *testing.T, db core_postgres.DB) {
	t.Helper()

	_, err := db.Exec(t.Context(), `
		TRUNCATE TABLE
			users
		RESTART IDENTITY
		CASCADE;
	`)
	if err != nil {
		t.Fatalf("reset database: %v", err)
	}
}

func LoadData(t *testing.T, db core_postgres.DB) {
	t.Helper()
	ResetDB(t, db)
	query := `
	INSERT INTO users (id, username, first_name, last_name, created_at, bio, password_hash)
	VALUES ($1, $2,$3,$4,$5,$6, $7) 
	`
	for _, user := range Users {
		_, err := db.Exec(
			t.Context(),
			query,
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName,
			CreatedAt,
			user.Bio,
			PasswordHash,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func GetLoggerContext(ctx context.Context) context.Context {
	return core_logger.WithLogger(ctx, Log)
}
