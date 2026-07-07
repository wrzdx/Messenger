package core_postgres

import "context"

type Pool interface {
	DB
	Begin(ctx context.Context) (Tx, error)
	Close()
}
