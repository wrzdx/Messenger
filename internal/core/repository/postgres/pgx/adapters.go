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

type pgxBatch struct {
	batch *pgx.Batch
}

func (b *pgxBatch) Queue(sql string, args ...any) {
	b.batch.Queue(sql, args...)
}

type pgxBatchResults struct {
	results pgx.BatchResults
}

func (br pgxBatchResults) Exec() (postgres.CommandTag, error) {
	tag, err := br.results.Exec()
	if err != nil {
		return nil, mapErrors(err)
	}
	return pgxCommandTag{tag}, nil
}

func (br pgxBatchResults) Query() (postgres.Rows, error) {
	rows, err := br.results.Query()
	if err != nil {
		return nil, mapErrors(err)
	}
	return pgxRows{rows}, nil
}

func (br pgxBatchResults) QueryRow() postgres.Row {
	return pgxRow{br.results.QueryRow()}
}

func (br pgxBatchResults) Close() error {
	return mapErrors(br.results.Close())
}

var violationErrs = map[string]error{
	"23503": postgres.ErrViolatesForeignKey,
	"23505": postgres.ErrViolatesUnique,
	"23514": postgres.ErrViolatesCheck,
	"22001": postgres.ErrTooLongVarchar,
}

func mapErrors(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return postgres.ErrNoRows
	}

	mappedErr := postgres.ErrUnknown
	var constraintName string

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		violationErr, ok := violationErrs[pgErr.Code]
		if ok {
			mappedErr = violationErr
			constraintName = pgErr.ConstraintName
		}
	}

	return postgres.DBError{
		Err:        mappedErr,
		Constraint: constraintName,
		Wrapped:    err,
	}
}
