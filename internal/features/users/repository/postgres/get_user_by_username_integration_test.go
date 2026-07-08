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

func TestGetUserByUsername(t *testing.T) {
	pool, err := pgx_pool.NewPool(t.Context(), pgx_pool.NewConfigMust())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer pool.Close()

	tests := []struct {
		name      string
		username  string
		before    func(t *testing.T, repo *UsersRepository)
		wantUser  domain.User
		wantError error
	}{
		{
			name:     "existing user",
			username: "ivanov",
			before: func(t *testing.T, repo *UsersRepository) {
				_, err := repo.CreateUser(
					t.Context(),
					domain.NewUser(
						uuid.New(),
						"ivanov",
						"Ivan",
						new("Ivanov"),
						test_utils.CreatedAt,
						new("I like pizza"),
						test_utils.PasswordHash,
					),
				)
				if err != nil {
					t.Fatal(err)
				}
			},
			wantUser: domain.NewUser(
				uuid.Nil, // сравним позже
				"ivanov",
				"Ivan",
				new("Ivanov"),
				test_utils.CreatedAt,
				new("I like pizza"),
				test_utils.PasswordHash,
			),
		},
		{
			name:      "user not found",
			username:  "unknown",
			wantError: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := pool.Begin(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(t.Context())

			repo := NewUsersRepository(tx)
			test_utils.LoadData(t, tx)

			if tt.before != nil {
				tt.before(t, repo)
			}

			gotUser, err := repo.GetUserByUsername(
				t.Context(),
				tt.username,
			)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("want %v, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := tt.wantUser
			want.ID = gotUser.ID

			if diff := cmp.Diff(want, gotUser); diff != "" {
				t.Fatalf("GetUserByUsername mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
