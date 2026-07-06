package users_service

import (
	"errors"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetUser(t *testing.T) {
	user := core_test_utils.Users[0]

	tests := []struct {
		name           string
		userID         int
		repoUser       domain.User
		repoErr        error
		wantUser       domain.User
		wantRepoID     int
		wantRepoCalled bool
		wantError      error
	}{
		{
			name:           "existing user",
			userID:         user.ID,
			repoUser:       user,
			wantUser:       user,
			wantRepoID:     user.ID,
			wantRepoCalled: true,
		},
		{
			name:           "non-existing user",
			userID:         -1,
			repoErr:        core_errors.ErrorNotFound,
			wantRepoID:     -1,
			wantRepoCalled: true,
			wantError:      core_errors.ErrorNotFound,
		},

		{
			name:           "repository error",
			userID:         user.ID,
			wantRepoID:     user.ID,
			wantRepoCalled: true,
			repoErr:        core_errors.ErrInternalServer,
			wantError:      core_errors.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				repoCalled bool
				repoGotID  int
			)

			repo := StubUsersRepository{
				GetUserFn: func(id int) (domain.User, error) {
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

			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
			}

			if diff := cmp.Diff(tt.wantUser, gotUser); diff != "" {
				t.Fatalf("GetUser mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
