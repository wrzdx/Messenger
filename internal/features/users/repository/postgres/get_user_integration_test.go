//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
	core_pgx_pool "messenger/internal/core/repository/postgres/pool/pgx"
	core_test_utils "messenger/internal/core/utils/test"
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
			userID:   core_test_utils.Users[0].ID,
			wantUser: core_test_utils.Users[0],
		},
		{
			name:      "non-existing user",
			userID:    core_test_utils.ID,
			wantError: core_postgres_pool.ErrNoRows,
		},
	}

	pool, err := core_pgx_pool.NewPool(t.Context(), core_pgx_pool.NewConfigMust())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer pool.Close()

	core_test_utils.LoadData(t, pool)

	repository := NewUsersRepository(pool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
