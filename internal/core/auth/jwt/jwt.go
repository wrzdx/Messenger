package core_auth_jwt

import (
	"fmt"
	core_auth "messenger/internal/core/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTProvider struct {
	config Config
}

func NewJWTProvider(config Config) *JWTProvider {
	return &JWTProvider{
		config: config,
	}
}

func (j *JWTProvider) generate(
	userID uuid.UUID,
	ttl time.Duration,
	tokenType core_auth.TokenType,
) (string, error) {
	now := time.Now()
	expires := now.Add(ttl)

	claims := jwtClaims{
		Claims: core_auth.Claims{
			UserID:    userID,
			Type:      tokenType,
			IssuedAt:  now,
			ExpiresAt: expires,
		},

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", fmt.Errorf(
			"sign token: %w",
			err,
		)
	}

	return signed, nil
}

func (j *JWTProvider) GenerateTokens(id uuid.UUID) (core_auth.AuthTokens, error) {
	access, err := j.generate(
		id,
		j.config.AccessTokenTTL,
		core_auth.TokenTypeAccess,
	)
	if err != nil {
		return core_auth.AuthTokens{}, err
	}
	refresh, err := j.generate(
		id,
		j.config.RefreshTokenTTL,
		core_auth.TokenTypeRefresh,
	)
	if err != nil {
		return core_auth.AuthTokens{}, err
	}
	return core_auth.AuthTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (j *JWTProvider) ParseToken(token string) (core_auth.Claims, error) {
	var claims jwtClaims

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (any, error) {
			return []byte(j.config.Secret), nil
		},
	)
	if err != nil {
		return core_auth.Claims{}, fmt.Errorf("parse JWT: %w", err)
	}

	return claims.Claims, nil
}
