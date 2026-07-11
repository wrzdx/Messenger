package postgres

import (
	"context"
	"fmt"
)

type txKey struct{}

type TransactionManager struct {
	pool Pool
}

func NewTransactionManager(pool Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

func (tm *TransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func GetExecutor(ctx context.Context, defaultPool DB) DB {
	if tx, ok := ctx.Value(txKey{}).(Tx); ok {
		return tx
	}
	return defaultPool
}
