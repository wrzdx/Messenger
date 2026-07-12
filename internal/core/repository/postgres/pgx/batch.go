package pgx_pool

import (
	postgres "messenger/internal/core/repository/postgres"

	"github.com/jackc/pgx/v5"
)

type BatchFactory struct{}

func (f BatchFactory) NewBatch() postgres.Batch {
	return &pgxBatch{batch: &pgx.Batch{}}
}
