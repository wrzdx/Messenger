package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaxUsernameLength  int = 32
	MinUsernameLength  int = 5
	MaxFirstnameLength int = 64
	MinFirstnameLength int = 1
	MaxLastNameLength  int = 64
	MaxBioLength       int = 70
	MaxPasswordLength  int = 32
	MinPasswordLength  int = 8
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
	fields := make(map[string]string)
	if err := ValidateUsername(u.Username); err != nil {
		fields["username"] = err.Error()
	}

	if err := ValidateFirstName(u.FirstName); err != nil {
		fields["first_name"] = err.Error()
	}

	if u.LastName != nil {
		if err := ValidateLastName(*u.LastName); err != nil {
			fields["last_name"] = err.Error()
		}
	}

	if u.Bio != nil {
		if err := ValidateBio(*u.Bio); err != nil {
			fields["bio"] = err.Error()
		}
	}
	if len(fields) > 0 {
		return ValidationErr(string(UserEntity), fields)
	}
	return nil
}

func ValidateUsername(username string) error {
	return validateLength(
		"username",
		username,
		new(MinUsernameLength),
		new(MaxUsernameLength),
	)
}
func ValidatePassword(password string) error {
	return validateLength(
		"password",
		password,
		new(MinPasswordLength),
		new(MaxPasswordLength),
	)
}
func ValidateFirstName(firstName string) error {
	return validateLength(
		"first_name",
		firstName,
		new(MinFirstnameLength),
		new(MaxFirstnameLength),
	)
}
func ValidateLastName(lastName string) error {
	return validateLength(
		"last_name",
		lastName,
		nil,
		new(MaxLastNameLength),
	)
}
func ValidateBio(bio string) error {
	return validateLength(
		"bio",
		bio,
		nil,
		new(MaxBioLength),
	)
}
