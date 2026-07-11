package auth_service

import (
	"context"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	test_utils "messenger/internal/core/utils/test"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	mockUser := test_utils.MockUser
	validTokens := auth.TokenPair{Access: "access", Refresh: "refresh"}
	tests := []struct {
		name       string
		payload    RegisterPayload
		setupMocks func(repo *MockUsersRepository, hasher *MockHasher, tokens *MockTokenProvider)
		wantUser   domain.User
		wantTokens auth.TokenPair
		wantError  error
	}{
		{
			name: "Success",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			setupMocks: func(repo *MockUsersRepository, hasher *MockHasher, tokens *MockTokenProvider) {
				hasher.EXPECT().Hash("fsociety").Return(mockUser.PasswordHash, nil)
				tokens.EXPECT().
					GenerateTokenPair(
						mock.AnythingOfType("uuid.UUID"),
						mock.AnythingOfType("uuid.UUID"),
					).Return(validTokens, nil)
				repo.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, user domain.User) (domain.User, error) {
						user.ID = mockUser.ID
						user.CreatedAt = mockUser.CreatedAt
						require.Equal(t, mockUser, user)
						return user, nil
					},
				)
			},
			wantTokens: validTokens,
			wantUser:   mockUser,
		},
		{
			name: "missing username",
			payload: RegisterPayload{
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "short username",
			payload: RegisterPayload{
				Username:  "ecor",
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "long username",
			payload: RegisterPayload{
				Username:  "ecorp" + strings.Repeat("a", 32),
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "missing firstname",
			payload: RegisterPayload{
				Username: "ecorp",
				LastName: new("Wellick"),
				Bio:      new("Dead"),
				Password: "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "long firstname",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell" + strings.Repeat("a", 64),
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "long lastname",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell",
				LastName:  new("Wellick" + strings.Repeat("a", 64)),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "long bio",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead" + strings.Repeat("a", 70)),
				Password:  "fsociety",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "short password",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociet",
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "long password",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety" + strings.Repeat("a", 32),
			},
			wantError: domain.ErrValidation,
		},
		{
			name: "Success",
			payload: RegisterPayload{
				Username:  "ecorp",
				FirstName: "Tyrell",
				LastName:  new("Wellick"),
				Bio:       new("Dead"),
				Password:  "fsociety",
			},
			setupMocks: func(repo *MockUsersRepository, hasher *MockHasher, tokens *MockTokenProvider) {
				hasher.EXPECT().Hash("fsociety").Return(mockUser.PasswordHash, nil)
				tokens.EXPECT().
					GenerateTokenPair(
						mock.AnythingOfType("uuid.UUID"),
						mock.AnythingOfType("uuid.UUID"),
					).Return(validTokens, nil)
				repo.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, user domain.User) (domain.User, error) {
						user.ID = mockUser.ID
						user.CreatedAt = mockUser.CreatedAt
						require.Equal(t, mockUser, user)
						return domain.User{}, domain.ErrAlreadyExists
					},
				)
			},
			wantError: domain.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUsersRepository(t)
			hasher := NewMockHasher(t)
			tokenProvider := NewMockTokenProvider(t)
			if tt.setupMocks != nil {
				tt.setupMocks(repo, hasher, tokenProvider)
			}

			service := NewAuthService(repo, hasher, tokenProvider)

			gotUser, gotTokens, err := service.Register(
				t.Context(),
				tt.payload,
			)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, gotUser)
			assert.Equal(t, tt.wantTokens, gotTokens)
		})
	}
}
