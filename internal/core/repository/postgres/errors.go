package postgres

import (
	"errors"
	"fmt"
)

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrViolatesUnique     = errors.New("violates unique constraint")
	ErrUnknown            = errors.New("unknown")
	ErrViolatesCheck      = errors.New("violates check constraint")
	ErrTooLongVarchar     = errors.New("value too long for type character varying")
)

type DBError struct {
	Err        error
	Constraint string
	Wrapped    error
}

func (e DBError) Error() string {
	return fmt.Sprintf("%v: %v", e.Err, e.Wrapped)
}

func (e DBError) Unwrap() error {
	return e.Err
}
