package domain

import "errors"

var (
	ErrNullUsername     = errors.New("username cannot be null")
	ErrInvalidUsername  = errors.New("username must be between 5 and 32 characters")
	ErrNullFirstname    = errors.New("first name cannot be null")
	ErrInvalidFirstName = errors.New("first name must be between 1 and 64 characters")
	ErrInvalidLastName  = errors.New("last name cannot exceed 64 characters")
	ErrInvalidBio       = errors.New("bio cannot exceed 70 characters")
	ErrInvalidPassword  = errors.New("password must be between 8 and 32 characters")

	ErrUserAlreadyExists  = errors.New("already exists")
	ErrUserNotFound       = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("wrong password")

	ErrNegativeLimit  = errors.New("limit must be non-negative")
	ErrNegativeOffset = errors.New("offset must be non-negative")
)
