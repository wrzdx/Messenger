package auth_service

// import (
// 	"errors"
// 	"messenger/internal/core/domain"
// 	test_utils "messenger/internal/core/utils/test"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/google/uuid"
// )

// func TestLogin(t *testing.T) {
// 	tests := []struct {
// 		name string

// 		username string
// 		password string

// 		user domain.User

// 		repoError   error
// 		hasherError error
// 		jwtError    error

// 		wantTokens domain.TokenPair
// 		wantError  error

// 		wantRepoCalled   bool
// 		wantHasherCalled bool
// 		wantJWTCalled    bool
// 	}{
// 		{
// 			name:     "valid credentials",
// 			username: "ivanov",
// 			password: "password",

// 			user: domain.NewUser(
// 				test_utils.ID,
// 				"ivanov",
// 				"Ivan",
// 				new("Ivanov"),
// 				test_utils.CreatedAt,
// 				new("bio"),
// 				test_utils.PasswordHash,
// 			),

// 			wantTokens: domain.TokenPair{
// 				Access:  "access-token",
// 				Refresh: "refresh-token",
// 			},

// 			wantRepoCalled:   true,
// 			wantHasherCalled: true,
// 			wantJWTCalled:    true,
// 		},
// 		{
// 			name:     "user not found",
// 			username: "ivanov",
// 			password: "password",

// 			repoError: domain.ErrUserNotFound,
// 			wantError: domain.ErrInvalidCredentials,

// 			wantRepoCalled:   true,
// 			wantHasherCalled: false,
// 			wantJWTCalled:    false,
// 		},
// 		{
// 			name:     "invalid password",
// 			username: "ivanov",
// 			password: "password",

// 			user: domain.NewUser(
// 				test_utils.ID,
// 				"ivanov",
// 				"Ivan",
// 				nil,
// 				test_utils.CreatedAt,
// 				nil,
// 				test_utils.PasswordHash,
// 			),

// 			hasherError: test_utils.HasherError,
// 			wantError:   test_utils.HasherError,

// 			wantRepoCalled:   true,
// 			wantHasherCalled: true,
// 			wantJWTCalled:    false,
// 		},
// 		{
// 			name:     "jwt error",
// 			username: "ivanov",
// 			password: "password",

// 			user: domain.NewUser(
// 				test_utils.ID,
// 				"ivanov",
// 				"Ivan",
// 				nil,
// 				test_utils.CreatedAt,
// 				nil,
// 				test_utils.PasswordHash,
// 			),

// 			jwtError:  test_utils.JWTError,
// 			wantError: test_utils.JWTError,

// 			wantRepoCalled:   true,
// 			wantHasherCalled: true,
// 			wantJWTCalled:    true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var (
// 				repoCalled   bool
// 				hasherCalled bool
// 				jwtCalled    bool

// 				repoUsername   string
// 				hasherHash     string
// 				hasherPassword string
// 				jwtUserID      = test_utils.ID
// 			)

// 			stubRepo := StubsUserRepository{
// 				GetUserByUsernameFn: func(
// 					username string,
// 				) (domain.User, error) {
// 					repoCalled = true
// 					repoUsername = username
// 					return tt.user, tt.repoError
// 				},
// 			}

// 			stubHasher := StubHasher{
// 				CompareFn: func(hash, password string) error {
// 					hasherCalled = true
// 					hasherHash = hash
// 					hasherPassword = password
// 					return tt.hasherError
// 				},
// 			}

// 			stubJWT := StubJWTProvider{
// 				GenerateTokensFn: func(id uuid.UUID) (domain.TokenPair, error) {
// 					jwtCalled = true
// 					jwtUserID = id
// 					return tt.wantTokens, tt.jwtError
// 				},
// 			}

// 			service := NewAuthService(&stubRepo, &stubHasher, &stubJWT)

// 			gotTokens, err := service.Login(
// 				t.Context(),
// 				tt.username,
// 				tt.password,
// 			)

// 			if repoCalled != tt.wantRepoCalled {
// 				t.Fatalf("repo called = %v, want %v", repoCalled, tt.wantRepoCalled)
// 			}

// 			if hasherCalled != tt.wantHasherCalled {
// 				t.Fatalf("hasher called = %v, want %v", hasherCalled, tt.wantHasherCalled)
// 			}

// 			if jwtCalled != tt.wantJWTCalled {
// 				t.Fatalf("jwt called = %v, want %v", jwtCalled, tt.wantJWTCalled)
// 			}

// 			if repoCalled && repoUsername != tt.username {
// 				t.Fatalf("repo username = %q, want %q", repoUsername, tt.username)
// 			}

// 			if hasherCalled {
// 				if hasherHash != tt.user.PasswordHash {
// 					t.Fatalf("hash = %q, want %q", hasherHash, tt.user.PasswordHash)
// 				}
// 				if hasherPassword != tt.password {
// 					t.Fatalf("password = %q, want %q", hasherPassword, tt.password)
// 				}
// 			}

// 			if jwtCalled && jwtUserID != tt.user.ID {
// 				t.Fatalf("jwt user id = %v, want %v", jwtUserID, tt.user.ID)
// 			}

// 			if tt.wantError != nil {
// 				if !errors.Is(err, tt.wantError) {
// 					t.Fatalf("want %v, got %v", tt.wantError, err)
// 				}
// 			} else if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}

// 			if diff := cmp.Diff(tt.wantTokens, gotTokens); diff != "" {
// 				t.Fatalf("tokens mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
