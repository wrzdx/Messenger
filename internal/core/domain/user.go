package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidUser = errors.New("invalid user")

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
	emptyProfile := UserProfile{}
	if u.Profile == emptyProfile {
		fields["profile"] = "profile cannot be empty"
	}

	if len(fields) > 0 {
		return DetailedError{
			Err:     ErrInvalidUser,
			Details: fields,
		}
	}
	return nil
}
