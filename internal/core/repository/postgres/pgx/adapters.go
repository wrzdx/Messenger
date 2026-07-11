package pgx_pool

import (
	"errors"
	postgres "messenger/internal/core/repository/postgres"

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
	"23503": postgres.ErrViolatesForeignKey,
	"23505": postgres.ErrViolatesUnique,
	"23514": postgres.ErrViolatesCheck,
	"22001": postgres.ErrTooLongVarchar,
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return postgres.ErrNoRows
	}

	mappedErr := postgres.ErrUnknown

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		violationErr, ok := violationErrs[pgErr.Code]
		if ok {
			mappedErr = violationErr
		}
	}

	return mappedErr
}
