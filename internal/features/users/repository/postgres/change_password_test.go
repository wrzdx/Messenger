//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	core_pgx_pool "messenger/internal/core/repository/postgres/pool/pgx"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/uuid"
)

func TestChangePassword(t *testing.T) {
	const newPasswordHash = "new_password_hash"

	tests := []struct {
		name         string
		userID       uuid.UUID
		passwordHash string
		wantError    error
	}{
		{
			name:         "existing user",
			userID:       core_test_utils.Users[0].ID,
			passwordHash: newPasswordHash,
		},
		{
			name:         "non-existing user",
			userID:       core_test_utils.ID,
			passwordHash: newPasswordHash,
			wantError:    domain.ErrUserNotFound,
		},
	}

	pool, err := core_pgx_pool.NewPool(
		t.Context(),
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer pool.Close()

	repository := NewUsersRepository(pool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core_test_utils.LoadData(t, pool)

			err := repository.ChangePassword(
				t.Context(),
				tt.userID,
				tt.passwordHash,
			)

			if !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}

			if tt.wantError == nil {
				user, err := repository.GetUser(t.Context(), tt.userID)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if user.PasswordHash != newPasswordHash {
					t.Fatalf(
						"password hash = %q, want %q",
						user.PasswordHash,
						newPasswordHash,
					)
				}
			}
		})
	}
}
