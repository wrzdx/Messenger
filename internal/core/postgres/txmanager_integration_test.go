//go:build integration

package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestTransactionManager(t *testing.T) {
	config := NewConfigMust()
	pool, err := NewPool(t.Context(), config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	tableName := pgx.Identifier{
		"transaction_manager_test_" + uuid.NewString(),
	}.Sanitize()
	_, err = pool.Exec(t.Context(), fmt.Sprintf(`
		CREATE TABLE %s (
			id INTEGER PRIMARY KEY
		)
	`, tableName))
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()
		_, err := pool.Exec(ctx, fmt.Sprintf(`DROP TABLE %s`, tableName))
		require.NoError(t, err)
	})

	manager := NewTransactionManager(pool)

	t.Run("commits when callback succeeds", func(t *testing.T) {
		err := manager.WithinTransaction(t.Context(), func(ctx context.Context) error {
			executor := GetExecutor(ctx, pool)
			_, isTransaction := executor.(pgx.Tx)
			require.True(t, isTransaction)

			_, err := executor.Exec(
				ctx,
				fmt.Sprintf(`INSERT INTO %s (id) VALUES ($1)`, tableName),
				1,
			)
			return err
		})

		require.NoError(t, err)
		require.Equal(t, 1, transactionManagerTestRowCount(t, pool, tableName, 1))
	})

	t.Run("rolls back when callback fails", func(t *testing.T) {
		callbackErr := errors.New("callback failure")
		err := manager.WithinTransaction(t.Context(), func(ctx context.Context) error {
			executor := GetExecutor(ctx, pool)
			_, err := executor.Exec(
				ctx,
				fmt.Sprintf(`INSERT INTO %s (id) VALUES ($1)`, tableName),
				2,
			)
			require.NoError(t, err)
			return callbackErr
		})

		require.ErrorIs(t, err, callbackErr)
		require.Equal(t, 0, transactionManagerTestRowCount(t, pool, tableName, 2))
	})
}

func transactionManagerTestRowCount(
	t *testing.T,
	db DBTX,
	tableName string,
	id int,
) int {
	t.Helper()

	var count int
	err := db.QueryRow(
		t.Context(),
		fmt.Sprintf(`SELECT count(*) FROM %s WHERE id = $1`, tableName),
		id,
	).Scan(&count)
	require.NoError(t, err)
	return count
}
