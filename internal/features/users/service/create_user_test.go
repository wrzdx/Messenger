package users_service

import (
	"errors"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tests = []struct {
	name string

	user        domain.User
	credentials domain.UserCredentials

	hasherError error
	repoError   error

	wantUser  domain.User
	wantError error

	wantHasherCalled bool
	wantRepoCalled   bool
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
		credentials: domain.NewCredentials(
			"ivanov",
			"password",
		),

		wantUser: domain.NewUser(
			core_test_utils.ID,
			"ivanov",
			"Ivan",
			new("Ivanov"),
			core_test_utils.CreatedAt,
			new("I like pizza"),
		),

		wantHasherCalled: true,
		wantRepoCalled:   true,
	},
	{
		name: "invalid user",

		user: domain.User{
			ID:        domain.UninitializedID,
			Username:  "abc", // слишком короткий
			FirstName: "Ivan",
			CreatedAt: core_test_utils.CreatedAt,
		},
		credentials: domain.NewCredentials(
			"abc",
			"password",
		),

		wantError: core_errors.ErrInvalidArgument,

		wantHasherCalled: false,
		wantRepoCalled:   false,
	},
	{
		name: "invalid credentials",

		user: domain.User{
			ID:        domain.UninitializedID,
			Username:  "ivanov",
			FirstName: "Ivan",
			CreatedAt: core_test_utils.CreatedAt,
		},
		credentials: domain.UserCredentials{
			Username: "ivanov",
			Password: "123", // слишком короткий
		},

		wantError: core_errors.ErrInvalidArgument,

		wantHasherCalled: false,
		wantRepoCalled:   false,
	},
	{
		name: "hasher error",

		user: domain.User{
			ID:        domain.UninitializedID,
			Username:  "ivanov",
			FirstName: "Ivan",
			CreatedAt: core_test_utils.CreatedAt,
		},
		credentials: domain.NewCredentials(
			"ivanov",
			"password",
		),

		hasherError: core_test_utils.HasherError,
		wantError:   core_test_utils.HasherError,

		wantHasherCalled: true,
		wantRepoCalled:   false,
	},
	{
		name: "repository error",

		user: domain.User{
			ID:        domain.UninitializedID,
			Username:  "ivanov",
			FirstName: "Ivan",
			CreatedAt: core_test_utils.CreatedAt,
		},
		credentials: domain.NewCredentials(
			"ivanov",
			"password",
		),

		repoError: core_test_utils.RepoError,
		wantError: core_test_utils.RepoError,

		wantHasherCalled: true,
		wantRepoCalled:   true,
	},
}

func TestCreateUser(t *testing.T) {
	for _, tt := range tests {
		// Setup
		t.Run(tt.name, func(t *testing.T) {
			var want domain.User
			if tt.wantError == nil {
				want = domain.NewUser(
					core_test_utils.ID,
					tt.user.Username,
					tt.user.FirstName,
					tt.user.LastName,
					core_test_utils.CreatedAt,
					tt.user.Bio,
				)
			}
			stubRepo := StubUsersRepository{
				ReturnUser:  want,
				ReturnError: tt.repoError,
			}

			pswHash := tt.credentials.Password + "_hash"

			stubHasher := StubHasher{
				ReturnHash:  []byte(pswHash),
				ReturnError: tt.hasherError,
			}

			wantRepoGotUser := domain.NewUser(
				domain.UninitializedID,
				tt.user.Username,
				tt.user.FirstName,
				tt.user.LastName,
				core_test_utils.CreatedAt,
				tt.user.Bio,
			)
			wantHasherGotPsw := tt.credentials.Password
			service := NewUsersService(&stubRepo, &stubHasher)
			ctx := t.Context()
			// Action
			gotUser, gotError := service.CreateUser(ctx, tt.user, tt.credentials)

			// Check
			if stubHasher.Called != tt.wantHasherCalled {
				t.Fatalf("Hasher called = %v, want %v",
					stubHasher.Called,
					tt.wantHasherCalled,
				)
			}
			if stubRepo.Called != tt.wantRepoCalled {
				t.Fatalf("Repository called = %v, want %v",
					stubRepo.Called,
					tt.wantRepoCalled,
				)
			}
			if stubHasher.Called {
				if diff := cmp.Diff(stubHasher.GotPassword, wantHasherGotPsw); diff != "" {
					t.Fatalf("HasherGotPassword mismatch (-want +got):\n%s", diff)
				}
			}
			if stubRepo.Called {
				if diff := cmp.Diff(stubRepo.GotUser, wantRepoGotUser); diff != "" {
					t.Fatalf("RepoGotUser mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(stubRepo.GotPswHash, pswHash); diff != "" {
					t.Fatalf("RepoGotPswHash mismatch (-want +got):\n%s", diff)
				}
			}
			if tt.wantError != nil {
				if !errors.Is(gotError, tt.wantError) {
					t.Fatalf("want %v, got %v", tt.wantError, gotError)
				}
			} else if gotError != nil {
				t.Fatalf("unexpected error: %v", gotError)
			}
			if diff := cmp.Diff(want, gotUser); diff != "" {
				t.Fatalf("CreateUser mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
