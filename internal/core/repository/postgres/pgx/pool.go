package pgx_pool

import (
	"context"
	"fmt"
	postgres "messenger/internal/core/repository/postgres"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	optTimout time.Duration
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		Pool:      pool,
		optTimout: config.Timeout,
	}, nil
}

func (p *Pool) OptTimeout() time.Duration {
	return p.optTimout
}

func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (postgres.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) postgres.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (postgres.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}

func (p *Pool) Begin(ctx context.Context) (postgres.Tx, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{
		tx:        tx,
		optTimout: p.optTimout,
	}, nil
}

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
