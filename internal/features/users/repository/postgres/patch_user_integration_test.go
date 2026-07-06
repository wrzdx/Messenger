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
			id:   core_test_utils.Users[0].ID,
			user: domain.NewUser(
				core_test_utils.Users[0].ID,
				"new_username",
				"NewName",
				new("NewLastName"),
				core_test_utils.Users[0].CreatedAt,
				new("New bio"),
				core_test_utils.PasswordHash,
			),
		},
		{
			name: "patch nullable fields to null",
			id:   core_test_utils.Users[0].ID,
			user: domain.NewUser(
				core_test_utils.Users[0].ID,
				core_test_utils.Users[0].Username,
				core_test_utils.Users[0].FirstName,
				nil,
				core_test_utils.Users[0].CreatedAt,
				nil,
				core_test_utils.PasswordHash,
			),
		},
		{
			name:      "user not found",
			id:        uuid.New(),
			user:      core_test_utils.Users[0],
			wantError: core_postgres_pool.ErrNoRows,
		},
		{
			name: "duplicate username",
			id:   core_test_utils.Users[0].ID,
			user: domain.NewUser(
				core_test_utils.Users[0].ID,
				"duplicate",
				"Ivan",
				nil,
				core_test_utils.Users[0].CreatedAt,
				nil,
				core_test_utils.PasswordHash,
			),
			wantError: core_postgres_pool.ErrViolatesUnique,
			before: func(t *testing.T, repo *UsersRepository) {
				_, err := repo.CreateUser(
					t.Context(),
					domain.NewUser(
						uuid.New(),
						"duplicate",
						"Ivan",
						nil,
						core_test_utils.CreatedAt,
						nil,
						core_test_utils.PasswordHash,
					),
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

	core_test_utils.LoadData(t, pool)
	repository := NewUsersRepository(pool)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core_test_utils.LoadData(t, pool)

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

			if diff := cmp.Diff(tt.user, gotUser); diff != "" {
				t.Fatalf("PatchUser mismatch (-want +got):\n%s", diff)
			}

			var userModel UserModel

			err := pool.QueryRow(
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

			if diff := cmp.Diff(tt.user, savedUser); diff != "" {
				t.Fatalf("saved user mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
