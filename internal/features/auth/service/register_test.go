package auth_service

// import (
// 	"errors"
// 	"messenger/internal/core/domain"
// 	test_utils "messenger/internal/core/utils/test"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// )

// var tests = []struct {
// 	name string

// 	user domain.RegisterUserPayload

// 	hasherError error
// 	repoError   error

// 	wantUser  domain.User
// 	wantError error

// 	wantHasherCalled bool
// 	wantRepoCalled   bool
// }{
// 	{
// 		name: "valid user",

// 		user: domain.NewRegisterUserPayload(
// 			"ivanov",
// 			"Ivan",
// 			new("Ivanov"),
// 			new("I like pizza"),
// 			"password",
// 		),

// 		wantUser: domain.NewUser(
// 			test_utils.ID,
// 			"ivanov",
// 			"Ivan",
// 			new("Ivanov"),
// 			test_utils.CreatedAt,
// 			new("I like pizza"),
// 			test_utils.PasswordHash,
// 		),

// 		wantHasherCalled: true,
// 		wantRepoCalled:   true,
// 	},
// 	{
// 		name: "invalid username",

// 		user: domain.NewRegisterUserPayload(
// 			"abc",
// 			"Ivan",
// 			new("Ivanov"),
// 			new("I like pizza"),
// 			"password",
// 		),

// 		wantError: domain.ErrInvalidUsername,

// 		wantHasherCalled: false,
// 		wantRepoCalled:   false,
// 	},
// 	{
// 		name: "invalid password",

// 		user: domain.NewRegisterUserPayload(
// 			"ivanov",
// 			"Ivan",
// 			new("Ivanov"),
// 			new("I like pizza"),
// 			"passwo",
// 		),

// 		wantError: domain.ErrInvalidPassword,

// 		wantHasherCalled: false,
// 		wantRepoCalled:   false,
// 	},
// 	{
// 		name: "hasher error",

// 		user: domain.NewRegisterUserPayload(
// 			"ivanov",
// 			"Ivan",
// 			new("Ivanov"),
// 			new("I like pizza"),
// 			"password",
// 		),

// 		hasherError: test_utils.HasherError,
// 		wantError:   test_utils.HasherError,

// 		wantHasherCalled: true,
// 		wantRepoCalled:   false,
// 	},
// 	{
// 		name: "repository error",

// 		user: domain.NewRegisterUserPayload(
// 			"ivanov",
// 			"Ivan",
// 			new("Ivanov"),
// 			new("I like pizza"),
// 			"password",
// 		),

// 		repoError: test_utils.RepoError,
// 		wantError: test_utils.RepoError,

// 		wantHasherCalled: true,
// 		wantRepoCalled:   true,
// 	},
// }

// func TestCreateUser(t *testing.T) {
// 	for _, tt := range tests {
// 		// Setup
// 		t.Run(tt.name, func(t *testing.T) {
// 			var want domain.User
// 			pswHash := tt.user.Password + "_hash"

// 			if tt.wantError == nil {
// 				want = domain.NewUser(
// 					tt.wantUser.ID,
// 					tt.user.Username,
// 					tt.user.FirstName,
// 					tt.user.LastName,
// 					test_utils.CreatedAt,
// 					tt.user.Bio,
// 					pswHash,
// 				)
// 			}
// 			var repoCalled bool
// 			var repoGotUser domain.User
// 			stubRepo := StubsUserRepository{
// 				CreateUserFn: func(
// 					user domain.User,
// 				) (domain.User, error) {
// 					repoCalled = true
// 					repoGotUser = user
// 					repoGotUser.ID = tt.wantUser.ID
// 					repoGotUser.CreatedAt = tt.wantUser.CreatedAt
// 					return want, tt.repoError
// 				},
// 			}

// 			var hasherCalled bool
// 			var hasherGotPsw string
// 			stubHasher := StubHasher{
// 				HashFn: func(password string) (string, error) {
// 					hasherCalled = true
// 					hasherGotPsw = password
// 					return pswHash, tt.hasherError
// 				},
// 			}

// 			stubJWTProvider := StubJWTProvider{}

// 			wantRepoGotUser := domain.NewUser(
// 				tt.wantUser.ID,
// 				tt.user.Username,
// 				tt.user.FirstName,
// 				tt.user.LastName,
// 				tt.wantUser.CreatedAt,
// 				tt.user.Bio,
// 				pswHash,
// 			)
// 			wantHasherGotPsw := tt.user.Password
// 			service := NewAuthService(&stubRepo, &stubHasher, &stubJWTProvider)
// 			ctx := t.Context()
// 			// Action
// 			gotUser, _, gotError := service.Register(ctx, tt.user)

// 			// Check
// 			if hasherCalled != tt.wantHasherCalled {
// 				t.Fatalf("Hasher called = %v, want %v",
// 					hasherCalled,
// 					tt.wantHasherCalled,
// 				)
// 			}
// 			if repoCalled != tt.wantRepoCalled {
// 				t.Fatalf("Repository called = %v, want %v",
// 					repoCalled,
// 					tt.wantRepoCalled,
// 				)
// 			}
// 			if hasherCalled {
// 				if diff := cmp.Diff(wantHasherGotPsw, hasherGotPsw); diff != "" {
// 					t.Fatalf("HasherGotPassword mismatch (-want +got):\n%s", diff)
// 				}
// 			}
// 			if repoCalled {
// 				if diff := cmp.Diff(wantRepoGotUser, repoGotUser); diff != "" {
// 					t.Fatalf("RepoGotUser mismatch (-want +got):\n%s", diff)
// 				}
// 			}
// 			if tt.wantError != nil {
// 				if !errors.Is(gotError, tt.wantError) {
// 					t.Fatalf("want %v, got %v", tt.wantError, gotError)
// 				}
// 			} else if gotError != nil {
// 				t.Fatalf("unexpected error: %v", gotError)
// 			}
// 			if diff := cmp.Diff(want, gotUser); diff != "" {
// 				t.Fatalf("CreateUser mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
