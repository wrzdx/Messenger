//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
	core_pgx_pool "messenger/internal/core/repository/postgres/pool/pgx"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	users = []domain.User{
		{
			ID:        1,
			Username:  "user_1",
			FirstName: "Username",
			LastName:  new("1"),
			CreatedAt: core_test_utils.CreatedAt,
			Bio:       new("I'm user 1"),
		},
		{
			ID:        2,
			Username:  "user_2",
			FirstName: "Username",
			LastName:  new("2"),
			CreatedAt: core_test_utils.CreatedAt,
			Bio:       new("I'm user 2"),
		},
		{
			ID:        3,
			Username:  "user_3",
			FirstName: "Username",
			LastName:  new("3"),
			CreatedAt: core_test_utils.CreatedAt,
			Bio:       new("I'm user 3"),
		},
	}
)

func TestGetUsers(t *testing.T) {
	// common setup

	var tests = []struct {
		name      string
		users     []domain.User
		limit     *int
		offset    *int
		wantError error
	}{
		{
			name:  "all users",
			users: users,
		},
		{
			name:  "limit users",
			users: users[:1],
			limit: new(1),
		},
		{
			name:      "negative limit",
			users:     nil,
			limit:     new(-1),
			wantError: core_errors.ErrInvalidArgument,
		},
		{
			name:   "offset users",
			users:  users[1:],
			offset: new(1),
		},
		{
			name:      "negative offset",
			users:     nil,
			offset:     new(-1),
			wantError: core_errors.ErrInvalidArgument,
		},
		{
			name:   "limit offset users",
			users:  users[1:2],
			limit:  new(1),
			offset: new(1),
		},
		{
			name:  "empty users",
			users: []domain.User{},
			limit: new(0),
		},
	}
	pool, err := core_pgx_pool.NewPool(t.Context(), core_pgx_pool.NewConfigMust())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer pool.Close()

	loadData(t, pool)
	repository := NewUsersRepository(pool)

	// subtests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// action
			gotUsers, gotErr := repository.GetUsers(t.Context(), tt.limit, tt.offset)

			// assertion
			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
			}

			if diff := cmp.Diff(tt.users, gotUsers); diff != "" {
				t.Fatalf("GetUsers mismatch got users (-want +got):\n%s", diff)
			}
		})
	}
}

func loadData(t *testing.T, pool core_postgres_pool.Pool) {
	t.Helper()
	core_test_utils.ResetDB(t, pool)
	query := `
	INSERT INTO users (username, first_name, last_name, created_at, bio, password_hash)
	VALUES ($1, $2,$3,$4,$5,$6) 
	`
	for _, user := range users {
		_, err := pool.Exec(
			t.Context(),
			query,
			user.Username,
			user.FirstName,
			user.LastName,
			core_test_utils.CreatedAt,
			user.Bio,
			core_test_utils.PasswordHash,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
