package postgres

import "context"

type Tx interface {
	DB
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
