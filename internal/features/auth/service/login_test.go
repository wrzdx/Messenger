package auth_service

import (
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	test_utils "messenger/internal/core/utils/test"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	mockUser := test_utils.MockUser
	validTokens := auth.TokenPair{Access: "access", Refresh: "refresh"}
	tests := []struct {
		name       string
		username   string
		password   string
		setupMocks func(repo *MockUsersRepository, hasher *MockHasher, tokens *MockTokenProvider)
		wantTokens auth.TokenPair
		wantError  error
	}{
		{
			name:     "Success",
			username: "ecorp",
			password: "fsociety",
			setupMocks: func(repo *MockUsersRepository, hasher *MockHasher, tokens *MockTokenProvider) {
				repo.EXPECT().GetUserByUsername(mock.Anything, "ecorp").Return(mockUser, nil)
				hasher.EXPECT().Compare(mockUser.PasswordHash, "fsociety").Return(nil)
				tokens.EXPECT().
					GenerateTokenPair(
						mockUser.ID,
						mock.AnythingOfType("uuid.UUID"),
					).Return(validTokens, nil)
			},
			wantTokens: validTokens,
			wantError:  nil,
		},
		{
			name:     "Wrong Password",
			username: "ecorp",
			password: "ffscoiety",
			setupMocks: func(repo *MockUsersRepository, hasher *MockHasher, tokens *MockTokenProvider) {
				repo.EXPECT().GetUserByUsername(mock.Anything, "ecorp").Return(mockUser, nil)
				hasher.EXPECT().Compare(mockUser.PasswordHash, "ffscoiety").Return(auth.ErrPasswordMismatch)
			},
			wantTokens: auth.TokenPair{},
			wantError:  domain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUsersRepository(t)
			hasher := NewMockHasher(t)
			tokenProvider := NewMockTokenProvider(t)
			tt.setupMocks(repo, hasher, tokenProvider)

			service := NewAuthService(repo, hasher, tokenProvider)

			gotTokens, err := service.Login(
				t.Context(),
				tt.username,
				tt.password,
			)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.wantTokens, gotTokens)
		})
	}
}
