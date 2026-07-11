package domain

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type Entity string

const (
	UserEntity       Entity = "user"
	ChatEntity       Entity = "chat"
	MessageEntity    Entity = "message"
	PaginationEntity Entity = "pagination options"
)

var (
	// Generic

	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrValidation    = errors.New("failed to validate")
)

func NotFoundErr(entity Entity, field, value string) error {
	return fmt.Errorf("%s with %s='%s' %w", entity, field, value, ErrNotFound)
}

func AlreadyExistsErr(entity Entity, details map[string]string) DetailedError {
	de := DetailedError{
		Err:     fmt.Errorf("%s %w", entity, ErrAlreadyExists),
		Details: details,
	}
	return de
}

func ValidationErr(entity Entity, details map[string]string) DetailedError {
	de := DetailedError{
		Err:     fmt.Errorf("%w %s", ErrValidation, entity),
		Details: details,
	}
	return de
}

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

func validateLength(field, value string, minLength, maxLength *int) error {
	l := utf8.RuneCountInString(value)
	if minLength != nil && l < *minLength {
		return fmt.Errorf("'%s' length must be at least %d character", field, *minLength)
	}
	if maxLength != nil && l > *maxLength {
		return fmt.Errorf("'%s' length must be at most %d character", field, *maxLength)
	}

	return nil
}
