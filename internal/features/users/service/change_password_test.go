package users_service

import (
	"context"
	"errors"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_ChangePassword(t *testing.T) {
	ctx := context.Background()

	id := uuid.New()

	user := domain.User{
		ID:           id,
		PasswordHash: "old-hash",
	}

	tests := []struct {
		name     string
		prepare  func(*MockUsersRepository, *MockHasher)
		wantErr  error
		errorMsg string
	}{
		{
			name: "success",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				hasher.EXPECT().
					Compare("old-hash", "old-password").
					Return(nil).
					Once()

				hasher.EXPECT().
					Hash("new-password").
					Return("new-hash", nil).
					Once()

				repo.EXPECT().
					ChangePassword(ctx, id, "new-hash").
					Return(nil).
					Once()
			},
		},
		{
			name: "get user error",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(domain.User{}, errors.New("db error")).
					Once()
			},
			errorMsg: "get user",
		},
		{
			name: "wrong password",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				hasher.EXPECT().
					Compare("old-hash", "old-password").
					Return(auth.ErrPasswordMismatch).
					Once()
			},
			wantErr: domain.ErrWrongPassword,
		},
		{
			name: "compare error",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				hasher.EXPECT().
					Compare("old-hash", "old-password").
					Return(errors.New("bcrypt error")).
					Once()
			},
			errorMsg: "compare passwords",
		},
		{
			name: "invalid new password",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				hasher.EXPECT().
					Compare("old-hash", "old-password").
					Return(nil).
					Once()
			},
			errorMsg: "validate new password",
		},
		{
			name: "hash error",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				hasher.EXPECT().
					Compare("old-hash", "old-password").
					Return(nil).
					Once()

				hasher.EXPECT().
					Hash("new-password").
					Return("", errors.New("hash error")).
					Once()
			},
			errorMsg: "hash password",
		},
		{
			name: "repository error",
			prepare: func(repo *MockUsersRepository, hasher *MockHasher) {
				repo.EXPECT().
					GetUser(ctx, id).
					Return(user, nil).
					Once()

				hasher.EXPECT().
					Compare("old-hash", "old-password").
					Return(nil).
					Once()

				hasher.EXPECT().
					Hash("new-password").
					Return("new-hash", nil).
					Once()

				repo.EXPECT().
					ChangePassword(ctx, id, "new-hash").
					Return(errors.New("db error")).
					Once()
			},
			errorMsg: "change user password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUsersRepository(t)
			hasher := NewMockHasher(t)

			tt.prepare(repo, hasher)

			service := NewUsersService(repo, hasher)

			newPassword := "new-password"
			if tt.name == "invalid new password" {
				newPassword = "123"
			}

			err := service.ChangePassword(
				ctx,
				id,
				"old-password",
				newPassword,
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}
