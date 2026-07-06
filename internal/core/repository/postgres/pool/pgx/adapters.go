package core_pgx_pool

import (
	"errors"
	"fmt"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

func (r pgxRows) Err() error {
	err := r.Rows.Err()
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

var violationErrs = map[string]error{
	"23503": core_postgres_pool.ErrViolatesForeignKey,
	"23505": core_postgres_pool.ErrViolatesUnique,
	"23514": core_postgres_pool.ErrViolatesCheck,
	"22001": core_postgres_pool.ErrTooLongVarchar,
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	mappedErr := core_postgres_pool.ErrUnknown

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		violationErr, ok := violationErrs[pgErr.Code]
		if ok {
			mappedErr = violationErr
		}
	}

	return fmt.Errorf("%w: %v", mappedErr, err)
}
