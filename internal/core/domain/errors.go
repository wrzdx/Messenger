package domain

import (
	"errors"
	"maps"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
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

func (e DetailedError) Fields() map[string]string {
	return maps.Clone(e.Details)
}
