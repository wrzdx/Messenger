package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	now := time.Now()
	profile, err := NewUserProfile("Username_1", "First Name", nil, nil)
	require.NoError(t, err)

	t.Run("creates valid user", func(t *testing.T) {
		id := uuid.New()

		user, err := NewUser(id, profile, now, nil, "password_hash")

		require.NoError(t, err)
		require.Equal(t, User{
			ID:           id,
			Profile:      profile,
			CreatedAt:    now,
			PasswordHash: "password_hash",
		}, user)
	})

	t.Run("returns zero user when invalid", func(t *testing.T) {
		user, err := NewUser(uuid.Nil, profile, now, nil, "password_hash")

		require.ErrorIs(t, err, ErrInvalidUser)
		require.Zero(t, user)
	})
}

func TestUserValidate(t *testing.T) {
	now := time.Now()
	profile, err := NewUserProfile("Username_1", "First Name", nil, nil)
	require.NoError(t, err)

	validUser := func() User {
		return User{
			ID:           uuid.New(),
			Profile:      profile,
			CreatedAt:    now,
			PasswordHash: "password_hash",
		}
	}

	tests := []struct {
		name      string
		change    func(*User)
		wantError error
	}{
		{name: "valid active user", change: func(*User) {}},
		{
			name: "valid user deleted after creation",
			change: func(user *User) {
				user.DeletedAt = new(now.Add(time.Hour))
			},
		},
		{
			name: "deleted_at equal to created_at",
			change: func(user *User) {
				user.DeletedAt = new(now)
			},
		},
		{
			name: "nil id",
			change: func(user *User) {
				user.ID = uuid.Nil
			},
			wantError: ErrInvalidUser,
		},
		{
			name: "zero created_at",
			change: func(user *User) {
				user.CreatedAt = time.Time{}
			},
			wantError: ErrInvalidUser,
		},
		{
			name: "zero deleted_at",
			change: func(user *User) {
				user.DeletedAt = new(time.Time{})
			},
			wantError: ErrInvalidUser,
		},
		{
			name: "deleted_at before created_at",
			change: func(user *User) {
				user.DeletedAt = new(now.Add(-time.Hour))
			},
			wantError: ErrInvalidUser,
		},
		{
			name: "empty password hash",
			change: func(user *User) {
				user.PasswordHash = ""
			},
			wantError: ErrInvalidUser,
		},
		{
			name: "empty profile",
			change: func(user *User) {
				user.Profile = UserProfile{}
			},
			wantError: ErrInvalidUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := validUser()
			tt.change(&user)

			require.ErrorIs(t, user.Validate(), tt.wantError)
		})
	}
}
