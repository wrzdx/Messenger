package domain

import "errors"

var (
	ErrInvalidUsername  = errors.New("username must be between 5 and 32 characters")
	ErrInvalidFirstName = errors.New("first name must be between 1 and 64 characters")
	ErrInvalidLastName  = errors.New("last name cannot exceed 64 characters")
	ErrInvalidBio       = errors.New("bio cannot exceed 70 characters")
	ErrInvalidPassword = errors.New("password must be between 8 and 32 characters")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
