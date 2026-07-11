package postgres

import (
	"context"
	"time"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	OptTimeout() time.Duration
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}

type Pool interface {
	DB
	Begin(ctx context.Context) (Tx, error)
	Close()
}

type Tx interface {
	DB
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
