package domain

import (
	"fmt"
	core_errors "messenger/internal/core/errors"
	"time"
)

type User struct {
	ID        int
	Username  string
	FirstName string
	LastName  *string
	CreatedAt time.Time
	Bio       *string
}

func NewUser(
	id int,
	username string,
	firstName string,
	lastName *string,
	createdAt time.Time,
	bio *string,
) User {
	return User{
		ID:        id,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: createdAt,
		Bio:       bio,
	}
}

func NewUserUninitialized(
	username string,
	firstName string,
	lastName *string,
	bio *string,
) User {
	return NewUser(
		UninitializedID,
		username,
		firstName,
		lastName,
		time.Now(),
		bio,
	)
}

func (u *User) Validate() error {
	usernameLen := len([]rune(u.Username))
	if usernameLen < 5 || usernameLen > 32 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			usernameLen,
			core_errors.ErrInvalidArgument,
		)
	}
	firstNameLen := len([]rune(u.FirstName))
	if firstNameLen < 1 || firstNameLen > 64 {
		return fmt.Errorf(
			"invalid `FirstName` len: %d: %w",
			firstNameLen,
			core_errors.ErrInvalidArgument,
		)
	}
	if u.LastName != nil {
		lastNameLen := len([]rune(*u.LastName))
		if lastNameLen > 64 {
			return fmt.Errorf(
				"invalid `LastName` len: %d: %w",
				lastNameLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if u.Bio != nil {
		bioLen := len([]rune(*u.Bio))
		if bioLen > 70 {
			return fmt.Errorf(
				"invalid `Bio` len: %d: %w",
				bioLen,
				core_errors.ErrInvalidArgument,
			)
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
		return fmt.Errorf(
			"`Username` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.FirstName.Set && p.FirstName.Value == nil {
		return fmt.Errorf(
			"`FirstName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
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
