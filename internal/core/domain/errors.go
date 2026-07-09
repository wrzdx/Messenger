package domain

import (
	"errors"
	"fmt"
)

type Entity string

const (
	UserEntity       Entity = "user"
)

var (
	// Generic

	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")

	// Pure

	ErrNullUsername     = errors.New("username cannot be null")
	ErrInvalidUsername  = errors.New("username must be between 5 and 32 characters")
	ErrNullFirstname    = errors.New("first_name cannot be null")
	ErrInvalidFirstName = errors.New("first_name must be between 1 and 64 characters")
	ErrInvalidLastName  = errors.New("last_name cannot exceed 64 characters")
	ErrInvalidBio       = errors.New("bio cannot exceed 70 characters")
	ErrInvalidPassword  = errors.New("password must be between 8 and 32 characters")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("wrong password")

	ErrNegativeLimit  = errors.New("limit must be non-negative")
	ErrNegativeOffset = errors.New("offset must be non-negative")
)

func NotFoundErr(entity Entity, field, value string) error {
	return fmt.Errorf("%s with %s='%s' %w", entity, field, value, ErrNotFound)
}

func AlreadyExistsErr(entity Entity, field, value string) error {
	return fmt.Errorf("%s with %s='%s' %w", entity, field, value, ErrAlreadyExists)
}
