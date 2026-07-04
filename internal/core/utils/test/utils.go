package core_test_utils

import (
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
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
