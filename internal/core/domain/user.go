package domain

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	FirstName    string
	LastName     *string
	CreatedAt    time.Time
	Bio          *string
	PasswordHash string
}

func NewUser(
	id uuid.UUID,
	username string,
	firstName string,
	lastName *string,
	createdAt time.Time,
	bio *string,
	passwordHash string,
) User {
	return User{
		ID:           id,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		CreatedAt:    createdAt,
		Bio:          bio,
		PasswordHash: passwordHash,
	}
}

func (u *User) Validate() []error {
	var errs []error
	if err := ValidateUsername(u.Username); err != nil {
		errs = append(errs, err)
	}

	if err := ValidateFirstName(u.FirstName); err != nil {
		errs = append(errs, err)
	}

	if u.LastName != nil {
		if err := ValidateLastName(*u.LastName); err != nil {
			errs = append(errs, err)
		}
	}

	if u.Bio != nil {
		if err := ValidateBio(*u.Bio); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

type UserPatch struct {
	Username  Nullable[string]
	FirstName Nullable[string]
	LastName  Nullable[string]
	Bio       Nullable[string]
}

func NewUserPatch(
	username Nullable[string],
	firstName Nullable[string],
	lastName Nullable[string],
	bio Nullable[string],
) UserPatch {
	return UserPatch{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Bio:       bio,
	}
}

func (p *UserPatch) Validate() []error {
	var errs []error
	user := User{}
	if p.Username.Set {
		if p.Username.Value == nil {
			errs = append(errs, ErrNullUsername)
		}
		user.Username = *p.Username.Value
	}

	if p.FirstName.Set {
		if p.FirstName.Value == nil {
			errs = append(errs, ErrNullFirstname)
		}
		user.FirstName = *p.FirstName.Value
	}
	if p.LastName.Set {
		user.LastName = p.LastName.Value
	}

	if p.Bio.Set {
		user.Bio = p.Bio.Value
	}

	if errs = append(errs, user.Validate()...); errs != nil {
		return errs
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if errs := patch.Validate(); errs != nil {
		return fmt.Errorf("validate user patch: %w", errors.Join(errs...))
	}

	tmp := *u

	if patch.Username.Set {
		tmp.Username = *patch.Username.Value
	}

	if patch.FirstName.Set {
		tmp.FirstName = *patch.FirstName.Value
	}

	if patch.LastName.Set {
		tmp.LastName = patch.LastName.Value
	}

	if patch.Bio.Set {
		tmp.Bio = patch.Bio.Value
	}

	if errs := tmp.Validate(); errs != nil {
		return fmt.Errorf("validate patched user: %w", errors.Join(errs...))
	}

	*u = tmp
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
