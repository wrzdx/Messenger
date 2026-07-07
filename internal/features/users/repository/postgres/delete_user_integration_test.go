//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	core_pgx_pool "messenger/internal/core/repository/postgres/pgx"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/uuid"
)

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name      string
		userID    uuid.UUID
		wantError error
	}{
		{
			name:   "existing user",
			userID: core_test_utils.Users[0].ID,
		},
		{
			name:      "non-existing user",
			userID:    core_test_utils.ID,
			wantError: domain.ErrUserNotFound,
		},
	}

	pool, err := core_pgx_pool.NewPool(t.Context(), core_pgx_pool.NewConfigMust())
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
			core_test_utils.LoadData(t, tx)

			gotErr := repository.DeleteUser(t.Context(), tt.userID)

			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
			}

			if tt.wantError == nil {
				_, err := repository.GetUser(t.Context(), tt.userID)
				if !errors.Is(err, domain.ErrUserNotFound) {
					t.Fatalf("user was not deleted")
				}
			}
		})
	}
}
