package domain

import (
	"errors"
	"unicode/utf8"

	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUsername  = errors.New("username must be between 5 and 32 characters")
	ErrInvalidFirstName = errors.New("first_name must be between 1 and 64 characters")
	ErrInvalidLastName  = errors.New("last_name cannot exceed 64 characters")
	ErrInvalidBio       = errors.New("bio cannot exceed 70 characters")
	ErrInvalidPassword  = errors.New("password must be between 8 and 32 characters")
)

type User struct {
	ID           uuid.UUID
	Username     string
	FirstName    string
	LastName     *string
	CreatedAt    time.Time
	DeletedAt    *time.Time
	Bio          *string
	PasswordHash string
}

func NewUser(
	id uuid.UUID,
	username string,
	firstName string,
	lastName *string,
	createdAt time.Time,
	deletedAt *time.Time,
	bio *string,
	passwordHash string,
) User {
	return User{
		ID:           id,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		CreatedAt:    createdAt,
		DeletedAt:    deletedAt,
		Bio:          bio,
		PasswordHash: passwordHash,
	}
}

func (u *User) Validate() error {
	if err := ValidateUsername(u.Username); err != nil {
		return err
	}

	if err := ValidateFirstName(u.FirstName); err != nil {
		return err
	}

	if u.LastName != nil {
		if err := ValidateLastName(*u.LastName); err != nil {
			return err
		}
	}

	if u.Bio != nil {
		if err := ValidateBio(*u.Bio); err != nil {
			return err
		}
	}

	return nil
}

func ValidateUsername(username string) error {
	if l := utf8.RuneCountInString(username); l < 5 || l > 32 {
		return ErrInvalidUsername
	}
	return nil
}
func ValidatePassword(password string) error {
	if l := utf8.RuneCountInString(password); l < 8 || l > 32 {
		return ErrInvalidPassword
	}
	return nil
}
func ValidateFirstName(firstName string) error {
	if l := utf8.RuneCountInString(firstName); l < 1 || l > 64 {
		return ErrInvalidFirstName
	}
	return nil
}
func ValidateLastName(lastName string) error {
	if l := utf8.RuneCountInString(lastName); l > 64 {
		return ErrInvalidLastName
	}
	return nil
}
func ValidateBio(bio string) error {
	if l := utf8.RuneCountInString(bio); l > 70 {
		return ErrInvalidBio
	}
	return nil
}
