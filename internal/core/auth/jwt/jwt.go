package auth_jwt

import (
	"fmt"
	"messenger/internal/core/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type accessClaims struct {
	Type   auth.TokenType `json:"type"`
	UserID uuid.UUID      `json:"user_id"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	Type      auth.TokenType `json:"type"`
	SessionID uuid.UUID      `json:"token_id"`
	jwt.RegisteredClaims
}

type TokenProvider struct {
	config Config
}

func NewTokenProvider(config Config) *TokenProvider {
	return &TokenProvider{
		config: config,
	}
}

func (s *TokenProvider) GenerateTokenPair(userID, tokenID uuid.UUID) (auth.TokenPair, error) {
	aClaims := accessClaims{
		Type:   auth.TokenTypeAccess,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, aClaims).SignedString(s.config.Secret)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("access: %w", err)
	}

	rClaims := refreshClaims{
		Type:      auth.TokenTypeRefresh,
		SessionID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims).SignedString(s.config.Secret)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("refresh: %w", err)
	}

	return auth.TokenPair{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (s *TokenProvider) ParseAccessToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return s.config.Secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*accessClaims)
	if !ok ||
		claims.Type != auth.TokenTypeAccess ||
		claims.ExpiresAt == nil ||
		claims.ExpiresAt.Time.Before(time.Now()) ||
		claims.UserID == uuid.Nil {
		return uuid.UUID{}, auth.ErrInvalidToken
	}

	return claims.UserID, nil
}

func (s *TokenProvider) ParseRefreshToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &refreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return s.config.Secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*refreshClaims)
	if !ok ||
		claims.Type != auth.TokenTypeRefresh ||
		claims.ExpiresAt == nil ||
		claims.ExpiresAt.Time.Before(time.Now()) ||
		claims.SessionID == uuid.Nil {
		return uuid.UUID{}, auth.ErrInvalidToken
	}

	return claims.SessionID, nil
}
