package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	CheckViolation      = "23514"
)

func IsConstraintViolation(err error, sqlState string, constraintName string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}

	return pgErr.Code == sqlState && pgErr.ConstraintName == constraintName
}
