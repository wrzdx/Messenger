package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrViolatesUnique     = errors.New("violates unique constraint")
	ErrUnknown            = errors.New("unknown")
	ErrViolatesCheck      = errors.New("violates check constraint")
	ErrTooLongVarchar     = errors.New("value too long for type character varying")
	ErrNegativeLimit      = errors.New("limit must be non-negative")
	ErrNegativeOffset      = errors.New("offset must be non-negative")
)
