package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidUser = errors.New("invalid user")
var ErrAlreadyDeleted = errors.New("user already deleted")

type User struct {
	ID           uuid.UUID
	Profile      UserProfile
	CreatedAt    time.Time
	DeletedAt    *time.Time
	PasswordHash string
}

func NewUser(
	id uuid.UUID,
	profile UserProfile,
	createdAt time.Time,
	deletedAt *time.Time,
	passwordHash string,
) (User, error) {
	user := User{
		ID:           id,
		Profile:      profile,
		CreatedAt:    createdAt,
		DeletedAt:    deletedAt,
		PasswordHash: passwordHash,
	}
	if err := user.Validate(); err != nil {
		return User{}, err
	}
	return user, nil
}

func (u User) Validate() error {
	fields := make(map[string]string)

	if u.ID == uuid.Nil {
		fields["id"] = "id cannot be nil"
	}

	if u.CreatedAt.IsZero() {
		fields["created_at"] = "created_at cannot be zero"
	}
	if u.DeletedAt != nil {
		if u.DeletedAt.IsZero() {
			fields["deleted_at"] = "deleted_at cannot be zero"
		} else if u.DeletedAt.Before(u.CreatedAt) {
			fields["deleted_at"] = "deleted_at cannot be before created_at"
		}
	}
	if u.PasswordHash == "" {
		fields["password_hash"] = "password_hash cannot be empty"
	}

	if err := u.Profile.Validate(); err != nil {
		return err
	}

	if len(fields) > 0 {
		return DetailedError{
			Err:     ErrInvalidUser,
			Details: fields,
		}
	}
	return nil
}

func (u User) Delete(deletedAt time.Time) (User, error) {
	if u.DeletedAt != nil {
		return User{}, ErrAlreadyDeleted
	}
	suffix := strings.ReplaceAll(u.ID.String(), "-", "")[:16]
	profile, err := NewUserProfile(
		"deleted_"+suffix,
		"Deleted Account",
		nil,
		nil,
	)
	if err != nil {
		return User{}, err
	}
	u.Profile = profile
	u.DeletedAt = &deletedAt
	if err := u.Validate(); err != nil {
		return User{}, err
	}
	return u, nil
}
