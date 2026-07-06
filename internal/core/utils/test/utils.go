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

func GetLoggerContext(r *http.Request) context.Context{
	log := core_logger.NewTestLogger()
	ctx := core_logger.ToContext(r.Context(), log)
	return ctx
}
