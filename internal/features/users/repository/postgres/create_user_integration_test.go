//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_pgx_pool "messenger/internal/core/repository/postgres/pool/pgx"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type test struct {
	name      string
	user      domain.User
	pswHash   string
	wantError error
	before    func(t *testing.T, repo *UsersRepository, tt test)
}

var tests = []test{
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
		pswHash: "some_hash",
	},
	{
		name: "without last name",
		user: domain.User{
			ID:        domain.UninitializedID,
			Username:  "petrov",
			FirstName: "Petr",
			CreatedAt: core_test_utils.CreatedAt,
		},
		pswHash: "some_hash",
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
		pswHash: "some_hash",
	},
	{
		name: "duplicate username",
		user: domain.User{
			ID:        domain.UninitializedID,
			Username:  "duplicate",
			FirstName: "Ivan",
			CreatedAt: core_test_utils.CreatedAt,
		},
		pswHash:   "hash",
		wantError: core_errors.ErrConflict,
		before: func(t *testing.T, repo *UsersRepository, tt test) {
			_, err := repo.CreateUser(
				t.Context(),
				tt.user,
				tt.pswHash,
			)
			if err != nil {
				t.Fatal(err)
			}
		},
	},
}

func TestCreateUser(t *testing.T) {
	// common setup
	pool, err := core_pgx_pool.NewPool(t.Context(), core_pgx_pool.NewConfigMust())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repository := NewUsersRepository(pool)

	// subtests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			core_test_utils.ResetDB(t, pool)
			if tt.before != nil {
				tt.before(t, repository, tt)
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
