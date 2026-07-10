package users_service

import (
	"context"
	"errors"
	"testing"

	"messenger/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_GetUsers(t *testing.T) {
	ctx := context.Background()

	limit := 10
	offset := 0

	users := []domain.User{
		{
			Username: "alice",
		},
		{
			Username: "bob",
		},
	}

	tests := []struct {
		name      string
		limit     *int
		offset    *int
		prepare   func(*MockUsersRepository)
		wantUsers []domain.User
		wantErr   error
		errorMsg  string
	}{
		{
			name:   "success",
			limit:  &limit,
			offset: &offset,
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUsers(
						ctx,
						&limit,
						&offset,
					).
					Return(users, nil).
					Once()
			},
			wantUsers: users,
		},
		{
			name:    "negative limit",
			limit:   new(-1),
			offset:  &offset,
			wantErr: domain.ErrNegativeLimit,
		},
		{
			name:    "negative offset",
			limit:   &limit,
			offset:  new(-1),
			wantErr: domain.ErrNegativeOffset,
		},
		{
			name:   "repository error",
			limit:  &limit,
			offset: &offset,
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUsers(
						ctx,
						&limit,
						&offset,
					).
					Return(nil, errors.New("database error")).
					Once()
			},
			errorMsg: "get users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUsersRepository(t)

			if tt.prepare != nil {
				tt.prepare(repo)
			}

			service := NewUsersService(repo, nil)

			gotUsers, err := service.GetUsers(
				ctx,
				tt.limit,
				tt.offset,
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
			assert.Equal(t, tt.wantUsers, gotUsers)
		})
	}
}
