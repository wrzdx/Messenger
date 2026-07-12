package pgx_pool

import (
	"context"
	postgres "messenger/internal/core/repository/postgres"
	"time"

	"github.com/jackc/pgx/v5"
)

type Tx struct {
	tx        pgx.Tx
	optTimout time.Duration
}

func (t *Tx) OptTimeout() time.Duration {
	return t.optTimout
}

func (t *Tx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (postgres.Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (t *Tx) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) postgres.Row {
	row := t.tx.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

func (t *Tx) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (postgres.CommandTag, error) {
	tag, err := t.tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}

func (t *Tx) Commit(ctx context.Context) error {
	return mapErrors(t.tx.Commit(ctx))
}

func (t *Tx) Rollback(ctx context.Context) error {
	return mapErrors(t.tx.Rollback(ctx))
}
func (t *Tx) SendBatch(ctx context.Context, b postgres.Batch) postgres.BatchResults {
	concreteBatch := b.(*pgxBatch)

	results := t.tx.SendBatch(ctx, concreteBatch.batch)
	return pgxBatchResults{results: results}
}
