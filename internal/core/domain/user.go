package domain

import (
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

func (u *User) Validate() error {
	if l := utf8.RuneCountInString(u.Username); l < 5 || l > 32 {
		return ErrInvalidUsername
	}
	if l := utf8.RuneCountInString(u.FirstName); l < 1 || l > 64 {
		return ErrInvalidFirstName
	}

	if u.LastName != nil {
		if l := utf8.RuneCountInString(*u.LastName); l > 64 {
			return ErrInvalidLastName
		}
	}

	if u.Bio != nil {
		if l := utf8.RuneCountInString(*u.Bio); l > 70 {
			return ErrInvalidBio
		}
	}

	return nil
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

func (p *UserPatch) Validate() error {
	if p.Username.Set && p.Username.Value == nil {
		return ErrInvalidUsername
	}
	if p.FirstName.Set && p.FirstName.Value == nil {
		return ErrInvalidFirstName
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
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

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp
	return nil
}
