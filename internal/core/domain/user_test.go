package domain

import (
	"strings"
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

func TestUserDelete(t *testing.T) {
	createdAt := time.Now().UTC()
	lastName := "Anderson"
	bio := "Profile to anonymize"
	profile, err := NewUserProfile("Username_1", "Elliot", &lastName, &bio)
	require.NoError(t, err)
	user, err := NewUser(uuid.New(), profile, createdAt, nil, "password_hash")
	require.NoError(t, err)

	t.Run("anonymizes active user", func(t *testing.T) {
		deletedAt := createdAt.Add(time.Hour)

		deleted, err := user.Delete(deletedAt)

		require.NoError(t, err)
		require.Equal(t, user.ID, deleted.ID)
		require.True(t, strings.HasPrefix(deleted.Profile.Username(), "deleted_"))
		require.Len(t, deleted.Profile.Username(), 24)
		require.True(t, UsernamePattern.MatchString(deleted.Profile.Username()))
		require.Equal(t, "Deleted Account", deleted.Profile.FirstName())
		require.Nil(t, deleted.Profile.LastName())
		require.Nil(t, deleted.Profile.Bio())
		require.NotNil(t, deleted.DeletedAt)
		require.True(t, deletedAt.Equal(*deleted.DeletedAt))
		require.True(t, user.CreatedAt.Equal(deleted.CreatedAt))
		require.Equal(t, user.PasswordHash, deleted.PasswordHash)
		require.Nil(t, user.DeletedAt)
		require.Equal(t, profile, user.Profile)
	})

	t.Run("rejects deletion before creation", func(t *testing.T) {
		deleted, err := user.Delete(createdAt.Add(-time.Second))

		require.ErrorIs(t, err, ErrInvalidUser)
		require.Zero(t, deleted)
	})

	t.Run("rejects zero deletion time", func(t *testing.T) {
		deleted, err := user.Delete(time.Time{})

		require.ErrorIs(t, err, ErrInvalidUser)
		require.Zero(t, deleted)
	})

	t.Run("rejects already deleted user", func(t *testing.T) {
		firstDeletedAt := createdAt.Add(time.Hour)
		deleted, err := user.Delete(firstDeletedAt)
		require.NoError(t, err)

		second, err := deleted.Delete(firstDeletedAt.Add(time.Hour))

		require.ErrorIs(t, err, ErrAlreadyDeleted)
		require.Zero(t, second)
		require.True(t, firstDeletedAt.Equal(*deleted.DeletedAt))
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
