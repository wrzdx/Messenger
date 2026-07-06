package core_test_utils

import (
	"context"
	core_logger "messenger/internal/core/logger"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
	"net/http"
	"testing"
)

func ResetDB(t *testing.T, pool core_postgres_pool.Pool) {
	t.Helper()

	_, err := pool.Exec(t.Context(), `
		TRUNCATE TABLE
			users
		RESTART IDENTITY
		CASCADE;
	`)
	if err != nil {
		t.Fatalf("reset database: %v", err)
	}
}

func LoadData(t *testing.T, pool core_postgres_pool.Pool) {
	t.Helper()
	ResetDB(t, pool)
	query := `
	INSERT INTO users (username, first_name, last_name, created_at, bio, password_hash)
	VALUES ($1, $2,$3,$4,$5,$6) 
	`
	for _, user := range Users {
		_, err := pool.Exec(
			t.Context(),
			query,
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

func GetLoggerContext(r *http.Request) context.Context {
	ctx := core_logger.ToContext(r.Context(), Log)
	return ctx
}
