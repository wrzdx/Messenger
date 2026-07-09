//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	pgx_pool "messenger/internal/core/repository/postgres/pgx"
	test_utils "messenger/internal/core/utils/test"
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
			userID:       test_utils.MockUsers[0].ID,
			passwordHash: newPasswordHash,
		},
		{
			name:         "non-existing user",
			userID:       test_utils.MockUser.ID,
			passwordHash: newPasswordHash,
			wantError:    domain.ErrNotFound,
		},
	}

	pool, err := pgx_pool.NewPool(
		t.Context(),
		pgx_pool.NewConfigMust(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer pool.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := pool.Begin(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(t.Context())
			repository := NewUsersRepository(tx)
			test_utils.LoadData(t, tx)

			err = repository.ChangePassword(
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
