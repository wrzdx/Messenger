package domain

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrValidation    = errors.New("failed to validate")
)

type DetailedError struct {
	Err     error
	Details map[string]string
}

func (e DetailedError) Error() string {
	return e.Err.Error()
}

func (e DetailedError) Unwrap() error {
	return e.Err
}
