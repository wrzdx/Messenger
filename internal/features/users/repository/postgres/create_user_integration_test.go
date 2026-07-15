//go:build integration

package users_postgres_repository

// import (
// 	"errors"
// 	"messenger/internal/core/domain"
// 	postgres "messenger/internal/core/repository/postgres"
// 	pgx_pool "messenger/internal/core/repository/postgres/pgx"
// 	test_utils "messenger/internal/core/utils/test"
// 	"strings"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/google/uuid"
// )

// func TestCreateUser(t *testing.T) {
// 	// common setup

// 	var tests = []struct {
// 		name      string
// 		user      domain.User
// 		wantError error
// 		before    func(t *testing.T, repo *UsersRepository)
// 	}{
// 		{
// 			name: "valid user",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "ivanov",
// 				FirstName:    "Ivan",
// 				LastName:     new("Ivanov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				Bio:          new("I like pizza"),
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 		},
// 		{
// 			name: "without username",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				FirstName:    "Sidor",
// 				LastName:     new("Sidorov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrViolatesCheck,
// 		},
// 		{
// 			name: "short username",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "ivan",
// 				FirstName:    "Sidor",
// 				LastName:     new("Sidorov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrViolatesCheck,
// 		},
// 		{
// 			name: "long username",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "ivanov" + strings.Repeat("R", 32),
// 				FirstName:    "Sidor",
// 				LastName:     new("Sidorov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrTooLongVarchar,
// 		},
// 		{
// 			name: "without firstname",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "ivanov",
// 				LastName:     new("Sidorov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrViolatesCheck,
// 		},
// 		{
// 			name: "long firstname",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "ivanov",
// 				FirstName:    "Sido" + strings.Repeat("R", 64),
// 				LastName:     new("Sidorov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrTooLongVarchar,
// 		},
// 		{
// 			name: "without last name",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "petrov",
// 				FirstName:    "Petr",
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 		},
// 		{
// 			name: "long last name",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "petrov",
// 				FirstName:    "Petr",
// 				LastName:     new("Sidorov" + strings.Repeat("R", 64)),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrTooLongVarchar,
// 		},
// 		{
// 			name: "without bio",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "sidorov",
// 				FirstName:    "Sidor",
// 				LastName:     new("Sidorov"),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 		},
// 		{
// 			name: "long bio",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "sidorov",
// 				FirstName:    "Sidor",
// 				Bio:          new("Sidorov" + strings.Repeat("R", 70)),
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: postgres.ErrTooLongVarchar,
// 		},
// 		{
// 			name: "duplicate username",
// 			user: domain.User{
// 				ID:           test_utils.MockUser.ID,
// 				Username:     "duplicate",
// 				FirstName:    "Ivan",
// 				CreatedAt:    test_utils.MockUser.CreatedAt,
// 				PasswordHash: test_utils.MockUser.PasswordHash,
// 			},
// 			wantError: domain.ErrAlreadyExists,
// 			before: func(t *testing.T, repo *UsersRepository) {
// 				_, err := repo.CreateUser(
// 					t.Context(),
// 					domain.User{
// 						ID:           uuid.New(),
// 						Username:     "duplicate",
// 						FirstName:    "Ivan",
// 						CreatedAt:    test_utils.MockUser.CreatedAt,
// 						PasswordHash: test_utils.MockUser.PasswordHash,
// 					},
// 				)

// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			},
// 		},
// 	}
// 	pool, err := pgx_pool.NewPool(t.Context(), pgx_pool.NewConfigMust())
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	defer pool.Close()

// 	// subtests
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// setup
// 			tx, err := pool.Begin(t.Context())
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer tx.Rollback(t.Context())
// 			repository := NewUsersRepository(tx)
// 			test_utils.ResetDB(t, tx)
// 			if tt.before != nil {
// 				tt.before(t, repository)
// 			}

// 			// action
// 			gotUser, gotError := repository.CreateUser(t.Context(), tt.user)

// 			// assertion
// 			if tt.wantError != nil {
// 				if !errors.Is(gotError, tt.wantError) {
// 					t.Fatalf("want %v, got %v", tt.wantError, gotError)
// 				}
// 				return
// 			} else if gotError != nil {
// 				t.Fatalf("unexpected error: %v", gotError)
// 			}
// 			want := domain.NewUser(
// 				gotUser.ID,
// 				tt.user.Username,
// 				tt.user.FirstName,
// 				tt.user.LastName,
// 				gotUser.CreatedAt,
// 				tt.user.DeletedAt,
// 				tt.user.Bio,
// 				tt.user.PasswordHash,
// 			)

// 			if diff := cmp.Diff(want, gotUser); diff != "" {
// 				t.Fatalf("CreateUser mismatch got user (-want +got):\n%s", diff)
// 			}

// 			var userModel UserModel
// 			query := `
// 			SELECT id, username, first_name, last_name, created_at,deleted_at, bio, password_hash
// 			FROM users WHERE id=$1;`
// 			err = tx.QueryRow(t.Context(), query, gotUser.ID).Scan(
// 				&userModel.ID,
// 				&userModel.Username,
// 				&userModel.FirstName,
// 				&userModel.LastName,
// 				&userModel.CreatedAt,
// 				&userModel.DeletedAt,
// 				&userModel.Bio,
// 				&userModel.PasswordHash,
// 			)
// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}
// 			savedUser := domain.NewUser(
// 				userModel.ID,
// 				userModel.Username,
// 				userModel.FirstName,
// 				userModel.LastName,
// 				userModel.CreatedAt,
// 				userModel.DeletedAt,
// 				userModel.Bio,
// 				userModel.PasswordHash,
// 			)
// 			want.ID = userModel.ID
// 			if diff := cmp.Diff(want, savedUser); diff != "" {
// 				t.Fatalf("CreateUser mismatch saved user (-want +got):\n%s", diff)
// 			}

// 			if tt.user.PasswordHash != userModel.PasswordHash {
// 				t.Fatalf(
// 					"Password hash mismatch: \nwant: %s\ngot: %s",
// 					tt.user.PasswordHash,
// 					userModel.PasswordHash,
// 				)
// 			}
// 		})
// 	}
// }
