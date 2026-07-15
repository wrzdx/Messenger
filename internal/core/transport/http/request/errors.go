package http_request

import (
	"errors"
	"maps"
)

var ErrInvalidRequest = errors.New("invalid request")

type fieldError struct {
	fields map[string]string
}

func newFieldError(fields map[string]string) fieldError {
	return fieldError{fields: maps.Clone(fields)}
}

func (e fieldError) Error() string {
	return ErrInvalidRequest.Error()
}

func (e fieldError) Unwrap() error {
	return ErrInvalidRequest
}

func (e fieldError) Fields() map[string]string {
	return maps.Clone(e.fields)
}
