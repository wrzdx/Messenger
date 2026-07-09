package pgx_pool

import (
	"errors"
	"fmt"
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

type dbErrorWithConstraint struct {
	constraint string
	err        error
	wrapped    error
}

func (e dbErrorWithConstraint) Error() string {
	return fmt.Sprintf("%v: %v", e.wrapped, e.err)
}

func (e dbErrorWithConstraint) Unwrap() error {
	return e.err
}

func (e dbErrorWithConstraint) Constraint() string {
	return e.constraint
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
	var constraintName string

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		violationErr, ok := violationErrs[pgErr.Code]
		if ok {
			mappedErr = violationErr
			constraintName = pgErr.ConstraintName
		}
	}

	return dbErrorWithConstraint{
		constraint: constraintName,
		err:        mappedErr,
		wrapped:    err,
	}
}
