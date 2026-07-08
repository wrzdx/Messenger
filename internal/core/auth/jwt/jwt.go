package auth_jwt

import (
	"fmt"
	"messenger/internal/core/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type accessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	UserID  uuid.UUID `json:"user_id"`
	TokenID uuid.UUID `json:"token_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	config Config
}

func NewTokenService(config Config) *TokenService {
	return &TokenService{
		config: config,
	}
}

func (s *TokenService) GenerateTokenPair(user auth.AccessClaims, tokenID uuid.UUID) (auth.TokenPair, error) {
	aClaims := accessClaims{
		UserID: user.UserID,
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
		UserID:  user.UserID,
		TokenID: tokenID,
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

func (s *TokenService) ParseAccessToken(tokenStr string) (auth.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(token *jwt.Token) (any, error) {
		return s.config.Secret, nil
	})
	if err != nil || !token.Valid {
		return auth.AccessClaims{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*accessClaims)
	if !ok {
		return auth.AccessClaims{}, auth.ErrInvalidToken
	}

	return auth.AccessClaims{
		UserID: claims.UserID,
	}, nil
}

func (s *TokenService) ParseRefreshToken(tokenStr string) (auth.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &refreshClaims{}, func(token *jwt.Token) (any, error) {
		return s.config.Secret, nil
	})
	if err != nil || !token.Valid {
		return auth.RefreshClaims{}, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*refreshClaims)
	if !ok {
		return auth.RefreshClaims{}, auth.ErrInvalidToken
	}

	return auth.RefreshClaims{
		UserID:  claims.UserID,
		TokenID: claims.TokenID,
	}, nil
}
