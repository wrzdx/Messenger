package users_service

import (
	"errors"
	"messenger/internal/core/domain"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestChangePassword(t *testing.T) {
	user := core_test_utils.Users[0]

	tests := []struct {
		name string

		userID      uuid.UUID
		oldPassword string
		newPassword string

		getUserErr        error
		compareErr        error
		hashErr           error
		changePasswordErr error

		wantGetUserCalled        bool
		wantCompareCalled        bool
		wantHashCalled           bool
		wantChangePasswordCalled bool

		wantError error
	}{
		{
			name:        "success",
			userID:      user.ID,
			oldPassword: "old_password",
			newPassword: "new_password",

			wantGetUserCalled:        true,
			wantCompareCalled:        true,
			wantHashCalled:           true,
			wantChangePasswordCalled: true,
		},
		{
			name:        "user not found",
			userID:      user.ID,
			oldPassword: "old_password",
			newPassword: "new_password",

			getUserErr: domain.ErrUserNotFound,

			wantGetUserCalled: true,

			wantError: domain.ErrUserNotFound,
		},
		{
			name:        "invalid credentials",
			userID:      user.ID,
			oldPassword: "wrong_password",
			newPassword: "new_password",

			compareErr: domain.ErrInvalidCredentials,

			wantGetUserCalled: true,
			wantCompareCalled: true,

			wantError: domain.ErrInvalidCredentials,
		},
		{
			name:        "invalid new password",
			userID:      user.ID,
			oldPassword: "old_password",
			newPassword: "123",

			wantGetUserCalled: true,
			wantCompareCalled: true,

			wantError: domain.ErrInvalidPassword,
		},
		{
			name:        "hash error",
			userID:      user.ID,
			oldPassword: "old_password",
			newPassword: "new_password",

			hashErr: core_test_utils.HasherError,

			wantGetUserCalled: true,
			wantCompareCalled: true,
			wantHashCalled:    true,

			wantError: core_test_utils.HasherError,
		},
		{
			name:        "repository error",
			userID:      user.ID,
			oldPassword: "old_password",
			newPassword: "new_password",

			changePasswordErr: core_test_utils.RepoError,

			wantGetUserCalled:        true,
			wantCompareCalled:        true,
			wantHashCalled:           true,
			wantChangePasswordCalled: true,

			wantError: core_test_utils.RepoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				getUserCalled bool

				compareCalled bool
				compareHash   string
				comparePass   string

				hashCalled bool
				hashPass   string

				changeCalled bool
				changeID     uuid.UUID
				changeHash   string
			)

			repo := StubUsersRepository{
				GetUserFn: func(id uuid.UUID) (domain.User, error) {
					getUserCalled = true
					return user, tt.getUserErr
				},
				ChangePasswordFn: func(id uuid.UUID, hash string) error {
					changeCalled = true
					changeID = id
					changeHash = hash
					return tt.changePasswordErr
				},
			}

			hasher := StubHasher{
				CompareFn: func(hash string, password string) error {
					compareCalled = true
					compareHash = hash
					comparePass = password
					return tt.compareErr
				},
				HashFn: func(password string) (string, error) {
					hashCalled = true
					hashPass = password
					return core_test_utils.PasswordHash, tt.hashErr
				},
			}

			service := NewUsersService(&repo, &hasher)

			err := service.ChangePassword(
				t.Context(),
				tt.userID,
				tt.oldPassword,
				tt.newPassword,
			)

			if getUserCalled != tt.wantGetUserCalled {
				t.Fatalf(
					"GetUser called = %v, want %v",
					getUserCalled,
					tt.wantGetUserCalled,
				)
			}

			if compareCalled != tt.wantCompareCalled {
				t.Fatalf(
					"Compare called = %v, want %v",
					compareCalled,
					tt.wantCompareCalled,
				)
			}

			if hashCalled != tt.wantHashCalled {
				t.Fatalf(
					"Hash called = %v, want %v",
					hashCalled,
					tt.wantHashCalled,
				)
			}

			if changeCalled != tt.wantChangePasswordCalled {
				t.Fatalf(
					"ChangePassword called = %v, want %v",
					changeCalled,
					tt.wantChangePasswordCalled,
				)
			}

			if compareCalled {
				if diff := cmp.Diff(user.PasswordHash, compareHash); diff != "" {
					t.Fatalf("compare hash mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.oldPassword, comparePass); diff != "" {
					t.Fatalf("compare password mismatch (-want +got):\n%s", diff)
				}
			}

			if hashCalled {
				if diff := cmp.Diff(tt.newPassword, hashPass); diff != "" {
					t.Fatalf("hash password mismatch (-want +got):\n%s", diff)
				}
			}

			if changeCalled {
				if diff := cmp.Diff(tt.userID, changeID); diff != "" {
					t.Fatalf("userID mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(core_test_utils.PasswordHash, changeHash); diff != "" {
					t.Fatalf("password hash mismatch (-want +got):\n%s", diff)
				}
			}

			if !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}
