package core_auth_jwt

import (
	"fmt"
	core_auth "messenger/internal/core/auth"
	core_errors "messenger/internal/core/errors.go"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	userID int,
	ttl time.Duration,
	tokenType core_auth.TokenType,
) (string, time.Time, error) {
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
		return "", time.Time{}, fmt.Errorf(
			"sign token: %v: %w",
			err,
			core_errors.ErrUnauthorized,
		)
	}

	return signed, expires, nil
}

func (j *JWTProvider) GenerateAccessToken(id int) (string, time.Time, error) {
	token, expires, err := j.generate(id, j.config.AccessTokenTTL, core_auth.TokenTypeAccess)
	return token, expires, err
}

func (j *JWTProvider) GenerateRefreshToken(id int) (string, time.Time, error) {
	token, expires, err := j.generate(id, j.config.RefreshTokenTTL, core_auth.TokenTypeRefresh)
	return token, expires, err
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
