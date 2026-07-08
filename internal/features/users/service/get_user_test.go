package users_service

import (
	"errors"
	"messenger/internal/core/domain"
	test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestGetUser(t *testing.T) {
	user := test_utils.Users[0]

	tests := []struct {
		name           string
		userID         uuid.UUID
		repoUser       domain.User
		repoErr        error
		wantRepoID     uuid.UUID
		wantRepoCalled bool
		wantUser       domain.User
		wantError      error
	}{
		{
			name:           "existing user",
			userID:         user.ID,
			repoUser:       user,
			wantRepoID:     user.ID,
			wantRepoCalled: true,
			wantUser:       user,
		},
		{
			name:           "non-existing user",
			userID:         user.ID,
			repoErr:        domain.ErrUserNotFound,
			wantRepoID:     user.ID,
			wantRepoCalled: true,
			wantUser:       domain.User{},
			wantError:      domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				repoCalled bool
				repoGotID  uuid.UUID
			)

			repo := StubUsersRepository{
				GetUserFn: func(
					id uuid.UUID,
				) (domain.User, error) {
					repoCalled = true
					repoGotID = id

					return tt.repoUser, tt.repoErr
				},
			}

			hasher := StubHasher{}
			service := NewUsersService(&repo, &hasher)

			gotUser, gotErr := service.GetUser(t.Context(), tt.userID)

			if repoCalled != tt.wantRepoCalled {
				t.Fatalf(
					"repository called = %v, want %v",
					repoCalled,
					tt.wantRepoCalled,
				)
			}

			if repoCalled {
				if diff := cmp.Diff(tt.wantRepoID, repoGotID); diff != "" {
					t.Fatalf("userID mismatch (-want +got):\n%s", diff)
				}
			}

			if diff := cmp.Diff(tt.wantUser, gotUser); diff != "" {
				t.Fatalf("user mismatch (-want +got):\n%s", diff)
			}

			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
			}
		})
	}
}
