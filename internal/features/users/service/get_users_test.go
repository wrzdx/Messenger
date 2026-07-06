package users_service

import (
	"errors"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var users = []domain.User{
	{
		ID:        1,
		Username:  "user_1",
		FirstName: "Username",
		LastName:  new("1"),
		CreatedAt: core_test_utils.CreatedAt,
		Bio:       new("I'm user 1"),
	},
	{
		ID:        2,
		Username:  "user_2",
		FirstName: "Username",
		LastName:  new("2"),
		CreatedAt: core_test_utils.CreatedAt,
		Bio:       new("I'm user 2"),
	},
	{
		ID:        3,
		Username:  "user_3",
		FirstName: "Username",
		LastName:  new("3"),
		CreatedAt: core_test_utils.CreatedAt,
		Bio:       new("I'm user 3"),
	},
}

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
			repoUsers:      users,
			wantUsers:      users,
		},
		{
			name:           "limit users",
			limit:          new(1),
			repoUsers:      users[:1],
			wantRepoCalled: true,
			wantUsers:      users[:1],
		},
		{
			name:      "negative limit",
			limit:     new(-1),
			wantError: core_errors.ErrInvalidArgument,
		},
		{
			name:           "offset users",
			offset:         new(1),
			repoUsers:      users[1:],
			wantRepoCalled: true,
			wantUsers:      users[1:],
		},
		{
			name:      "negative offset",
			offset:     new(-1),
			wantError: core_errors.ErrInvalidArgument,
		},
		{
			name:           "limit offset users",
			limit:          new(1),
			offset:         new(1),
			repoUsers:      users[1:2],
			wantRepoCalled: true,

			wantUsers: users[1:2],
		},
		{
			name:           "empty users",
			limit:          new(1),
			offset:         new(2),
			repoUsers:      users[2:2],
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
			hasher := StubHasher{}
			service := NewUsersService(&repo, &hasher)

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
