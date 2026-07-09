package users_service

import (
	"context"
	"errors"
	"testing"

	"messenger/internal/core/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_PatchUser(t *testing.T) {
	ctx := context.Background()

	id := uuid.New()

	user := domain.User{
		ID:        id,
		Username:  "ecorp",
		FirstName: "Elliot",
	}

	updated := user
	updated.Username = "fsociety"

	username := "fsociety"

	validPatch := domain.NewUserPatch(
		domain.Nullable[string]{Value: &username, Set: true},
		domain.Nullable[string]{},
		domain.Nullable[string]{},
		domain.Nullable[string]{},
	)

	tests := []struct {
		name      string
		patch     domain.UserPatch
		prepare   func(*MockUsersRepository)
		wantUser  domain.User
		wantErr   error
		errorText string
	}{
		{
			name:  "success",
			patch: validPatch,
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				repo.EXPECT().
					PatchUser(ctx, id, updated).
					Return(updated, nil).
					Once()
			},
			wantUser: updated,
		},
		{
			name:  "get user error",
			patch: validPatch,
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(domain.User{}, errors.New("db")).Once()

			},
			errorText: "get user",
		},
		{
			name: "invalid username",
			patch: domain.NewUserPatch(
				domain.Nullable[string]{Value: new("abc"), Set: true},
				domain.Nullable[string]{},
				domain.Nullable[string]{},
				domain.Nullable[string]{},
			),
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()
			},
			wantErr: domain.ErrInvalidUsername,
		},
		{
			name:  "repository error",
			patch: validPatch,
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				repo.EXPECT().
					PatchUser(ctx, id, updated).
					Return(domain.User{}, errors.New("db")).Once()
			},
			errorText: "patch user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUsersRepository(t)

			tt.prepare(repo)

			service := NewUsersService(repo, nil)

			got, err := service.PatchUser(
				ctx,
				id,
				tt.patch,
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if tt.errorText != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorText)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, got)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
