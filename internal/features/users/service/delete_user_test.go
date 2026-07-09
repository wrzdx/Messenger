package users_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	tests := []struct {
		name     string
		prepare  func(*MockUsersRepository)
		wantErr  error
		errorMsg string
	}{
		{
			name: "success",
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					DeleteUser(ctx, id).
					Return(nil).
					Once()
			},
		},
		{
			name: "repository error",
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					DeleteUser(ctx, id).
					Return(errors.New("database error")).
					Once()
			},
			errorMsg: "delete user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUsersRepository(t)

			tt.prepare(repo)

			service := NewUsersService(repo, nil)

			err := service.DeleteUser(ctx, id)

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
