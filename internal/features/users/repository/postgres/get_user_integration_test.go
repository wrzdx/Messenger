//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	pgx_pool "messenger/internal/core/repository/postgres/pgx"
	test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name      string
		userID    uuid.UUID
		wantUser  domain.User
		wantError error
	}{
		{
			name:     "existing user",
			userID:   test_utils.Users[0].ID,
			wantUser: test_utils.Users[0],
		},
		{
			name:      "non-existing user",
			userID:    test_utils.ID,
			wantError: domain.ErrUserNotFound,
		},
	}

	pool, err := pgx_pool.NewPool(t.Context(), pgx_pool.NewConfigMust())
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
			test_utils.LoadData(t, tx)

			repository := NewUsersRepository(tx)
			gotUser, gotErr := repository.GetUser(t.Context(), tt.userID)
			tt.wantUser.PasswordHash = gotUser.PasswordHash
			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
			}

			if diff := cmp.Diff(tt.wantUser, gotUser); diff != "" {
				t.Fatalf("GetUser mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
