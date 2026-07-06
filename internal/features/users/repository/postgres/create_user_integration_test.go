//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_postgres_pool "messenger/internal/core/repository/postgres/pool"
	core_pgx_pool "messenger/internal/core/repository/postgres/pool/pgx"
	core_test_utils "messenger/internal/core/utils/test"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateUser(t *testing.T) {
	// common setup

	var tests = []struct {
		name      string
		user      domain.User
		pswHash   string
		wantError error
		before    func(t *testing.T, repo *UsersRepository)
	}{
		{
			name: "valid user",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "ivanov",
				FirstName: "Ivan",
				LastName:  new("Ivanov"),
				CreatedAt: core_test_utils.CreatedAt,
				Bio:       new("I like pizza"),
			},
			pswHash: core_test_utils.PasswordHash,
		},
		{
			name: "without username",
			user: domain.User{
				ID:        domain.UninitializedID,
				FirstName: "Sidor",
				LastName:  new("Sidorov"),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrViolatesCheck,
		},
		{
			name: "short username",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "ivan",
				FirstName: "Sidor",
				LastName:  new("Sidorov"),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrViolatesCheck,
		},
		{
			name: "long username",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "ivanov" + strings.Repeat("R", 32),
				FirstName: "Sidor",
				LastName:  new("Sidorov"),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrTooLongVarchar,
		},
		{
			name: "without firstname",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "ivanov",
				LastName:  new("Sidorov"),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrViolatesCheck,
		},
		{
			name: "long firstname",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "ivanov",
				FirstName: "Sido" + strings.Repeat("R", 64),
				LastName:  new("Sidorov"),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrTooLongVarchar,
		},
		{
			name: "without last name",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "petrov",
				FirstName: "Petr",
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
		},
		{
			name: "long last name",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "petrov",
				FirstName: "Petr",
				LastName:  new("Sidorov" + strings.Repeat("R", 64)),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrTooLongVarchar,
		},
		{
			name: "without bio",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "sidorov",
				FirstName: "Sidor",
				LastName:  new("Sidorov"),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
		},
		{
			name: "long bio",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "sidorov",
				FirstName: "Sidor",
				Bio:  new("Sidorov" + strings.Repeat("R", 70)),
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash: core_test_utils.PasswordHash,
			wantError: core_postgres_pool.ErrTooLongVarchar,
		},
		{
			name: "duplicate username",
			user: domain.User{
				ID:        domain.UninitializedID,
				Username:  "duplicate",
				FirstName: "Ivan",
				CreatedAt: core_test_utils.CreatedAt,
			},
			pswHash:   core_test_utils.PasswordHash,
			wantError: core_errors.ErrConflict,
			before: func(t *testing.T, repo *UsersRepository) {
				_, err := repo.CreateUser(
					t.Context(),
					domain.User{
						ID:        domain.UninitializedID,
						Username:  "duplicate",
						FirstName: "Ivan",
						CreatedAt: core_test_utils.CreatedAt,
					},
					core_test_utils.PasswordHash,
				)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	}
	pool, err := core_pgx_pool.NewPool(t.Context(), core_pgx_pool.NewConfigMust())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer pool.Close()
	repository := NewUsersRepository(pool)

	// subtests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			core_test_utils.ResetDB(t, pool)
			if tt.before != nil {
				tt.before(t, repository)
			}

			// action
			gotUser, gotError := repository.CreateUser(t.Context(), tt.user, tt.pswHash)

			// assertion
			if tt.wantError != nil {
				if !errors.Is(gotError, tt.wantError) {
					t.Fatalf("want %v, got %v", tt.wantError, gotError)
				}
				return
			} else if gotError != nil {
				t.Fatalf("unexpected error: %v", gotError)
			}
			want := domain.NewUser(
				gotUser.ID,
				tt.user.Username,
				tt.user.FirstName,
				tt.user.LastName,
				tt.user.CreatedAt,
				tt.user.Bio,
			)

			if diff := cmp.Diff(want, gotUser); diff != "" {
				t.Fatalf("CreateUser mismatch got user (-want +got):\n%s", diff)
			}

			var userModel UserModel
			query := `
			SELECT id, username, first_name, last_name, created_at, bio, password_hash 
			FROM users WHERE id=$1;`
			err := pool.QueryRow(t.Context(), query, gotUser.ID).Scan(
				&userModel.ID,
				&userModel.Username,
				&userModel.FirstName,
				&userModel.LastName,
				&userModel.CreatedAt,
				&userModel.Bio,
				&userModel.PasswordHash,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			savedUser := domain.NewUser(
				userModel.ID,
				userModel.Username,
				userModel.FirstName,
				userModel.LastName,
				userModel.CreatedAt,
				userModel.Bio,
			)
			want.ID = userModel.ID
			if diff := cmp.Diff(want, savedUser); diff != "" {
				t.Fatalf("CreateUser mismatch saved user (-want +got):\n%s", diff)
			}

			if tt.pswHash != userModel.PasswordHash {
				t.Fatalf(
					"Password hash mismatch: \nwant: %s\ngot: %s",
					tt.pswHash,
					userModel.PasswordHash,
				)
			}
		})
	}
}
