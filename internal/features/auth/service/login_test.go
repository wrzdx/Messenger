package auth_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Run("creates session and returns token pair", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		sessionsRepository := NewMockSessionsRepository(t)
		hasher := NewMockHasher(t)
		tokenProvider := NewMockTokenProvider(t)
		user := newLoginTestUser(t, false)
		config := AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SessionTTL:     24 * time.Hour,
		}
		var accessClaims auth.AccessTokenClaims
		var accessLifetime auth.TokenLifetime
		var refreshClaims auth.RefreshTokenClaims
		var refreshLifetime auth.TokenLifetime
		var savedSession domain.Session

		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().
			Compare(user.PasswordHash, "valid password value").
			Return(nil)
		tokenProvider.EXPECT().
			GenerateAccessToken(mock.Anything, mock.Anything).
			Run(func(claims auth.AccessTokenClaims, lifetime auth.TokenLifetime) {
				accessClaims = claims
				accessLifetime = lifetime
			}).
			Return("access-token", nil)
		tokenProvider.EXPECT().
			GenerateRefreshToken(mock.Anything, mock.Anything).
			Run(func(claims auth.RefreshTokenClaims, lifetime auth.TokenLifetime) {
				refreshClaims = claims
				refreshLifetime = lifetime
			}).
			Return("refresh-token", nil)
		sessionsRepository.EXPECT().
			CreateSession(mock.Anything, mock.Anything).
			Run(func(_ context.Context, session domain.Session) {
				savedSession = session
			}).
			Return(nil)
		service := &AuthService{
			usersRepository:    usersRepository,
			sessionsRepository: sessionsRepository,
			hasher:             hasher,
			tokenProvider:      tokenProvider,
			config:             config,
		}

		tokens, err := service.Login(t.Context(), "Username_1", "valid password value")

		require.NoError(t, err)
		require.Equal(t, auth.TokenPair{
			Access:  "access-token",
			Refresh: "refresh-token",
		}, tokens)
		require.NotEqual(t, uuid.Nil, savedSession.ID)
		require.Equal(t, user.ID, savedSession.UserID)
		require.NotEqual(t, uuid.Nil, savedSession.CurrentTokenID)
		require.True(t, savedSession.CreatedAt.Equal(savedSession.LastUsedAt))
		require.Equal(t, config.SessionTTL, savedSession.ExpiresAt.Sub(savedSession.CreatedAt))
		require.Equal(t, user.ID, accessClaims.UserID)
		require.True(t, savedSession.CreatedAt.Equal(accessLifetime.IssuedAt))
		require.Equal(t, config.AccessTokenTTL, accessLifetime.ExpiresAt.Sub(accessLifetime.IssuedAt))
		require.Equal(t, savedSession.ID, refreshClaims.SessionID)
		require.Equal(t, savedSession.CurrentTokenID, refreshClaims.TokenID)
		require.True(t, savedSession.CreatedAt.Equal(refreshLifetime.IssuedAt))
		require.True(t, savedSession.ExpiresAt.Equal(refreshLifetime.ExpiresAt))
	})

	t.Run("returns invalid credentials for unknown username", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "unknown").
			Return(domain.User{}, domain.ErrNotFound)
		hasher.EXPECT().DummyCompare().Return()
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
		}

		_, err := service.Login(t.Context(), "unknown", "password")

		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("returns user repository error without dummy comparison", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		repositoryErr := errors.New("database unavailable")
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(domain.User{}, repositoryErr)
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
		}

		_, err := service.Login(t.Context(), "Username_1", "password")

		require.ErrorIs(t, err, repositoryErr)
	})

	t.Run("returns invalid credentials for deleted user", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		user := newLoginTestUser(t, true)
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().DummyCompare().Return()
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
		}

		_, err := service.Login(t.Context(), "Username_1", "password")

		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("returns invalid credentials for password mismatch", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		user := newLoginTestUser(t, false)
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().
			Compare(user.PasswordHash, "wrong password").
			Return(auth.ErrPasswordMismatch)
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
		}

		_, err := service.Login(t.Context(), "Username_1", "wrong password")

		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("returns unexpected password comparison error", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		user := newLoginTestUser(t, false)
		compareErr := errors.New("password verifier failure")
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().
			Compare(user.PasswordHash, "password").
			Return(compareErr)
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
		}

		_, err := service.Login(t.Context(), "Username_1", "password")

		require.ErrorIs(t, err, compareErr)
	})

	t.Run("does not persist session when access token generation fails", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		tokenProvider := NewMockTokenProvider(t)
		user := newLoginTestUser(t, false)
		generateErr := errors.New("access signer failure")
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().Compare(user.PasswordHash, "password").Return(nil)
		tokenProvider.EXPECT().
			GenerateAccessToken(mock.Anything, mock.Anything).
			Return("", generateErr)
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
			tokenProvider:   tokenProvider,
			config: AuthConfig{
				AccessTokenTTL: 15 * time.Minute,
				SessionTTL:     24 * time.Hour,
			},
		}

		_, err := service.Login(t.Context(), "Username_1", "password")

		require.ErrorIs(t, err, generateErr)
	})

	t.Run("does not persist session when refresh token generation fails", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		hasher := NewMockHasher(t)
		tokenProvider := NewMockTokenProvider(t)
		user := newLoginTestUser(t, false)
		generateErr := errors.New("refresh signer failure")
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().Compare(user.PasswordHash, "password").Return(nil)
		tokenProvider.EXPECT().
			GenerateAccessToken(mock.Anything, mock.Anything).
			Return("access-token", nil)
		tokenProvider.EXPECT().
			GenerateRefreshToken(mock.Anything, mock.Anything).
			Return("", generateErr)
		service := &AuthService{
			usersRepository: usersRepository,
			hasher:          hasher,
			tokenProvider:   tokenProvider,
			config: AuthConfig{
				AccessTokenTTL: 15 * time.Minute,
				SessionTTL:     24 * time.Hour,
			},
		}

		_, err := service.Login(t.Context(), "Username_1", "password")

		require.ErrorIs(t, err, generateErr)
	})

	t.Run("returns session persistence error", func(t *testing.T) {
		usersRepository := NewMockUsersRepository(t)
		sessionsRepository := NewMockSessionsRepository(t)
		hasher := NewMockHasher(t)
		tokenProvider := NewMockTokenProvider(t)
		user := newLoginTestUser(t, false)
		createErr := errors.New("create session")
		usersRepository.EXPECT().
			GetUserByUsername(mock.Anything, "Username_1").
			Return(user, nil)
		hasher.EXPECT().Compare(user.PasswordHash, "password").Return(nil)
		tokenProvider.EXPECT().
			GenerateAccessToken(mock.Anything, mock.Anything).
			Return("access-token", nil)
		tokenProvider.EXPECT().
			GenerateRefreshToken(mock.Anything, mock.Anything).
			Return("refresh-token", nil)
		sessionsRepository.EXPECT().
			CreateSession(mock.Anything, mock.Anything).
			Return(createErr)
		service := &AuthService{
			usersRepository:    usersRepository,
			sessionsRepository: sessionsRepository,
			hasher:             hasher,
			tokenProvider:      tokenProvider,
			config: AuthConfig{
				AccessTokenTTL: 15 * time.Minute,
				SessionTTL:     24 * time.Hour,
			},
		}

		_, err := service.Login(t.Context(), "Username_1", "password")

		require.ErrorIs(t, err, createErr)
	})
}

func newLoginTestUser(t *testing.T, deleted bool) domain.User {
	t.Helper()

	profile, err := domain.NewUserProfile("Username_1", "First Name", nil, nil)
	require.NoError(t, err)
	createdAt := time.Date(2026, time.July, 15, 10, 0, 0, 0, time.UTC)
	var deletedAt *time.Time
	if deleted {
		deletedAt = new(createdAt.Add(time.Hour))
	}
	user, err := domain.NewUser(
		uuid.New(),
		profile,
		createdAt,
		deletedAt,
		"password-hash",
	)
	require.NoError(t, err)
	return user
}
