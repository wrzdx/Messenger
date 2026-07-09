package users_service

// import (
// 	"errors"
// 	"messenger/internal/core/domain"
// 	test_utils "messenger/internal/core/utils/test"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/google/uuid"
// )

// func TestDeleteUser(t *testing.T) {
// 	user := test_utils.Users[0]

// 	tests := []struct {
// 		name           string
// 		userID         uuid.UUID
// 		repoErr        error
// 		wantRepoID     uuid.UUID
// 		wantRepoCalled bool
// 		wantError      error
// 	}{
// 		{
// 			name:           "existing user",
// 			userID:         user.ID,
// 			wantRepoID:     user.ID,
// 			wantRepoCalled: true,
// 		},
// 		{
// 			name:           "non-existing user",
// 			repoErr:        domain.ErrUserNotFound,
// 			wantRepoCalled: true,
// 			wantError:      domain.ErrUserNotFound,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var (
// 				repoCalled bool
// 				repoGotID  uuid.UUID
// 			)

// 			repo := StubUsersRepository{
// 				DeleteUserFn: func(id uuid.UUID) error {
// 					repoCalled = true
// 					repoGotID = id

// 					return tt.repoErr
// 				},
// 			}

// 			hasher := StubHasher{}
// 			service := NewUsersService(&repo, &hasher)

// 			gotErr := service.DeleteUser(t.Context(), tt.userID)

// 			if repoCalled != tt.wantRepoCalled {
// 				t.Fatalf(
// 					"repository called = %v, want %v",
// 					repoCalled,
// 					tt.wantRepoCalled,
// 				)
// 			}

// 			if repoCalled {
// 				if diff := cmp.Diff(tt.wantRepoID, repoGotID); diff != "" {
// 					t.Fatalf("userID mismatch (-want +got):\n%s", diff)
// 				}
// 			}

// 			if !errors.Is(gotErr, tt.wantError) {
// 				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
// 			}
// 		})
// 	}
// }
