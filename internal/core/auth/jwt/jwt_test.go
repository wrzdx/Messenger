package auth_jwt

import (
	"messenger/internal/core/auth"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseAccessToken(t *testing.T) {
	secret := []byte("secret")

	tests := []struct {
		name   string
		claims accessClaims

		wantError error
	}{
		{
			name: "correct_token",
			claims: accessClaims{
				auth.TokenTypeAccess,
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
		},
		{
			name: "invalid type",
			claims: accessClaims{
				"invalid type",
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
		{
			name: "nil uuid",
			claims: accessClaims{
				auth.TokenTypeAccess,
				uuid.Nil,
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
		{
			name: "expired",
			claims: accessClaims{
				auth.TokenTypeAccess,
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
		{
			name: "nil expires at",
			claims: accessClaims{
				auth.TokenTypeAccess,
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: nil,
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims).SignedString(secret)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			tokenProvider := NewTokenProvider(Config{Secret: secret})
			userID, err := tokenProvider.ParseAccessToken(accessToken)
			require.ErrorIs(t, tt.wantError, err)
			if tt.wantError == nil {
				require.Equal(t, tt.claims.UserID, userID)
			}
		},
		)
	}

	t.Run("invalid claims type", func(t *testing.T) {
		claims := jwt.RegisteredClaims{
			ExpiresAt: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now())}

		accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tokenProvider := NewTokenProvider(Config{Secret: secret})
		_, err = tokenProvider.ParseAccessToken(accessToken)
		require.ErrorIs(t, auth.ErrInvalidToken, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		claims := accessClaims{
			auth.TokenTypeAccess,
			uuid.New(),
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now())},
		}

		accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		accessToken += "err"
		tokenProvider := NewTokenProvider(Config{Secret: secret})
		_, err = tokenProvider.ParseAccessToken(accessToken)
		require.ErrorIs(t, auth.ErrInvalidToken, err)
	})
}

func TestParseRefreshToken(t *testing.T) {
	secret := []byte("secret")

	tests := []struct {
		name   string
		claims refreshClaims

		wantError error
	}{
		{
			name: "correct_token",
			claims: refreshClaims{
				auth.TokenTypeRefresh,
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
		},
		{
			name: "invalid type",
			claims: refreshClaims{
				"invalid type",
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
		{
			name: "nil uuid",
			claims: refreshClaims{
				auth.TokenTypeRefresh,
				uuid.Nil,
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
		{
			name: "expired",
			claims: refreshClaims{
				auth.TokenTypeRefresh,
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
		{
			name: "nil expires at",
			claims: refreshClaims{
				auth.TokenTypeRefresh,
				uuid.New(),
				jwt.RegisteredClaims{
					ExpiresAt: nil,
					IssuedAt:  jwt.NewNumericDate(time.Now())},
			},
			wantError: auth.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims).SignedString(secret)
			require.NoError(t, err)

			tokenProvider := NewTokenProvider(Config{Secret: secret})
			tokenID, err := tokenProvider.ParseRefreshToken(refresh)
			require.ErrorIs(t, tt.wantError, err)
			if tt.wantError == nil {
				require.Equal(t, tt.claims.SessionID, tokenID)
			}
		},
		)
	}

	t.Run("invalid claims type", func(t *testing.T) {
		claims := jwt.RegisteredClaims{
			ExpiresAt: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now())}

		refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)

		tokenProvider := NewTokenProvider(Config{Secret: secret})
		_, err = tokenProvider.ParseRefreshToken(refresh)
		require.ErrorIs(t, auth.ErrInvalidToken, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		claims := refreshClaims{
			auth.TokenTypeRefresh,
			uuid.New(),
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now())},
		}

		refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)
		refresh += "err"
		tokenProvider := NewTokenProvider(Config{Secret: secret})
		_, err = tokenProvider.ParseRefreshToken(refresh)
		require.ErrorIs(t, auth.ErrInvalidToken, err)
	})

	t.Run("access token cannot be used as refresh token", func(t *testing.T) {
		tokenProvider := NewTokenProvider(Config{
			Secret:          secret,
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		})

		claims := accessClaims{
			auth.TokenTypeAccess,
			uuid.New(),
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now())},
		}
		access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)
		_, err = tokenProvider.ParseRefreshToken(access)
		require.ErrorIs(t, err, auth.ErrInvalidToken)
	})
	t.Run("refresh token cannot be used as access token", func(t *testing.T) {
		tokenProvider := NewTokenProvider(Config{
			Secret:          secret,
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		})

		claims := refreshClaims{
			auth.TokenTypeRefresh,
			uuid.New(),
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now())},
		}

		refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)

		_, err = tokenProvider.ParseAccessToken(refresh)
		require.ErrorIs(t, err, auth.ErrInvalidToken)
	})
}

func TestGenerateTokenPair(t *testing.T) {
	userID := uuid.New()
	tokenID := uuid.New()
	t.Run("correct", func(t *testing.T) {
		tokenProvider := NewTokenProvider(Config{
			Secret:          []byte("secret"),
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		})
		tokenPair, err := tokenProvider.GenerateTokenPair(userID, tokenID)
		require.NoError(t, err)
		gotUserID, err := tokenProvider.ParseAccessToken(tokenPair.Access)
		require.NoError(t, err)
		gotTokenID, err := tokenProvider.ParseRefreshToken(tokenPair.Refresh)
		require.NoError(t, err)
		require.Equal(t, userID, gotUserID)
		require.Equal(t, tokenID, gotTokenID)
	})
	t.Run("invalid access ttl", func(t *testing.T) {
		tokenProvider := NewTokenProvider(Config{
			Secret:         []byte("secret"),
			AccessTokenTTL: -time.Hour,
		})
		tokenPair, err := tokenProvider.GenerateTokenPair(userID, tokenID)
		require.NoError(t, err)
		_, err = tokenProvider.ParseAccessToken(tokenPair.Access)
		require.ErrorIs(t, auth.ErrInvalidToken, err)
	})
	t.Run("invalid refresh ttl", func(t *testing.T) {
		tokenProvider := NewTokenProvider(Config{
			Secret:          []byte("secret"),
			RefreshTokenTTL: -time.Hour,
		})
		tokenPair, err := tokenProvider.GenerateTokenPair(userID, tokenID)
		require.NoError(t, err)
		_, err = tokenProvider.ParseRefreshToken(tokenPair.Refresh)
		require.ErrorIs(t, auth.ErrInvalidToken, err)
	})
}
