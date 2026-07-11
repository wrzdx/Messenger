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
		name       string
		pagination domain.Pagination
		prepare    func(*MockUsersRepository)
		wantUsers  []domain.User
		wantErr    error
		errorMsg   string
	}{
		{
			name: "success",
			pagination: domain.NewPagination(
				&limit,
				&offset,
			),
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUsers(
						ctx,
						domain.NewPagination(&limit, &offset),
					).
					Return(users, nil).
					Once()
			},
			wantUsers: users,
		},
		{
			name: "negative limit",
			pagination: domain.NewPagination(
				new(-1),
				&offset,
			),
			wantErr: domain.ErrValidation,
		},
		{
			name: "negative offset",
			pagination: domain.NewPagination(
				&limit,
				new(-1),
			),
			wantErr: domain.ErrValidation,
		},
		{
			name: "repository error",
			pagination: domain.NewPagination(
				&limit,
				&offset,
			),
			prepare: func(repo *MockUsersRepository) {
				repo.EXPECT().
					GetUsers(
						ctx,
						domain.NewPagination(&limit, &offset),
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
				tt.pagination,
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
