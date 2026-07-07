package users_service

import (
	"errors"
	"messenger/internal/core/domain"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name           string
		limit          *int
		offset         *int
		repoUsers      []domain.User
		repoErr        error
		wantUsers      []domain.User
		wantRepoCalled bool
		wantError      error
	}{
		{
			name:           "return all users",
			wantRepoCalled: true,
			repoUsers:      core_test_utils.Users,
			wantUsers:      core_test_utils.Users,
		},
		{
			name:           "limit users",
			limit:          new(1),
			repoUsers:      core_test_utils.Users[:1],
			wantRepoCalled: true,
			wantUsers:      core_test_utils.Users[:1],
		},
		{
			name:           "offset users",
			offset:         new(1),
			repoUsers:      core_test_utils.Users[1:],
			wantRepoCalled: true,
			wantUsers:      core_test_utils.Users[1:],
		},
		{
			name:           "limit offset users",
			limit:          new(1),
			offset:         new(1),
			repoUsers:      core_test_utils.Users[1:2],
			wantRepoCalled: true,

			wantUsers: core_test_utils.Users[1:2],
		},
		{
			name:           "empty users",
			limit:          new(1),
			offset:         new(2),
			repoUsers:      core_test_utils.Users[2:2],
			wantRepoCalled: true,
			wantUsers:      []domain.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			var (
				repoCalled    bool
				repoGotLimit  *int
				repoGotOffset *int
			)

			repo := StubUsersRepository{
				GetUsersFn: func(limit, offset *int) ([]domain.User, error) {
					repoCalled = true
					repoGotLimit = limit
					repoGotOffset = offset

					return tt.repoUsers, tt.repoErr
				},
			}

			service := NewUsersService(&repo)

			// action
			gotUsers, gotErr := service.GetUsers(t.Context(), tt.limit, tt.offset)

			// check
			if repoCalled != tt.wantRepoCalled {
				t.Fatalf("Repository called = %v, want %v",
					repoCalled,
					tt.wantRepoCalled,
				)
			}

			if repoCalled {
				if diff := cmp.Diff(tt.limit, repoGotLimit); diff != "" {
					t.Fatalf("limit mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.offset, repoGotOffset); diff != "" {
					t.Fatalf("offset mismatch (-want +got):\n%s", diff)
				}
			}

			if !errors.Is(gotErr, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
			}
			if diff := cmp.Diff(tt.wantUsers, gotUsers); diff != "" {
				t.Fatalf("GetUsers mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
