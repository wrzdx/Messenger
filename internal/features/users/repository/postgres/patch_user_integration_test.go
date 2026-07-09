//go:build integration

package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
	postgres "messenger/internal/core/repository/postgres"
	pgx_pool "messenger/internal/core/repository/postgres/pgx"
	test_utils "messenger/internal/core/utils/test"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestPatchUser(t *testing.T) {
	var tests = []struct {
		name      string
		id        uuid.UUID
		user      domain.User
		wantError error
		before    func(t *testing.T, repo *UsersRepository)
	}{
		{
			name: "patch all fields",
			id:   test_utils.MockUsers[0].ID,
			user: domain.NewUser(
				test_utils.MockUsers[0].ID,
				"new_username",
				"NewName",
				new("NewLastName"),
				test_utils.MockUsers[0].CreatedAt,
				new("New bio"),
				test_utils.MockUsers[0].PasswordHash,
			),
		},
		{
			name: "patch nullable fields to null",
			id:   test_utils.MockUsers[0].ID,
			user: domain.NewUser(
				test_utils.MockUsers[0].ID,
				test_utils.MockUsers[0].Username,
				test_utils.MockUsers[0].FirstName,
				nil,
				test_utils.MockUsers[0].CreatedAt,
				nil,
				test_utils.MockUsers[0].PasswordHash,
			),
		},
		{
			name:      "user not found",
			id:        uuid.New(),
			wantError: postgres.ErrNoRows,
		},
		{
			name: "duplicate username",
			id:   test_utils.MockUsers[0].ID,
			user: domain.NewUser(
				test_utils.MockUsers[0].ID,
				"duplicate",
				"Ivan",
				nil,
				test_utils.MockUsers[0].CreatedAt,
				nil,
				test_utils.MockUsers[0].PasswordHash,
			),
			wantError: domain.ErrAlreadyExists,
			before: func(t *testing.T, repo *UsersRepository) {
				_, err := repo.CreateUser(
					t.Context(),
					domain.NewUser(
						uuid.New(),
						"duplicate",
						"Ivan",
						nil,
						test_utils.MockUsers[0].CreatedAt,
						nil,
						test_utils.MockUsers[0].PasswordHash,
					),
				)
				if err != nil {
					t.Fatal(err)
				}
			},
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
			test_utils.LoadData(t, tx)

			if tt.before != nil {
				tt.before(t, repository)
			}

			gotUser, gotErr := repository.PatchUser(
				t.Context(),
				tt.id,
				tt.user,
			)

			if tt.wantError != nil {
				if !errors.Is(gotErr, tt.wantError) {
					t.Fatalf("want %v, got %v", tt.wantError, gotErr)
				}
				return
			}

			if diff := cmp.Diff(tt.user, gotUser, cmpopts.EquateApproxTime(time.Microsecond)); diff != "" {
				t.Fatalf("PatchUser mismatch (-want +got):\n%s", diff)
			}

			var userModel UserModel

			err = tx.QueryRow(
				t.Context(),
				`
    SELECT id, username, first_name, last_name, created_at, bio, password_hash
    FROM users
    WHERE id=$1
    `,
				tt.id,
			).Scan(
				&userModel.ID,
				&userModel.Username,
				&userModel.FirstName,
				&userModel.LastName,
				&userModel.CreatedAt,
				&userModel.Bio,
				&userModel.PasswordHash,
			)
			if err != nil {
				t.Fatal(err)
			}

			savedUser := domain.NewUser(
				userModel.ID,
				userModel.Username,
				userModel.FirstName,
				userModel.LastName,
				userModel.CreatedAt,
				userModel.Bio,
				userModel.PasswordHash,
			)

			if diff := cmp.Diff(tt.user, savedUser, cmpopts.EquateApproxTime(time.Microsecond)); diff != "" {
				t.Fatalf("saved user mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
