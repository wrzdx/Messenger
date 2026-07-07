package core_auth_jwt

import (
	"fmt"
	"messenger/internal/core/domain"
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
	tokenType domain.TokenType,
) (domain.Token, error) {
	now := time.Now()
	expires := now.Add(ttl)

	claims := jwtClaims{
		UserID:    userID,
		Type:      tokenType,
		IssuedAt:  now,
		ExpiresAt: expires,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return domain.Token{}, fmt.Errorf(
			"sign token: %w",
			err,
		)
	}

	domainToken := domain.NewToken(signed, expires)

	return domainToken, nil
}

func (j *JWTProvider) GenerateAccessToken(id uuid.UUID) (domain.Token, error) {
	return j.generate(
		id,
		j.config.AccessTokenTTL,
		domain.TokenTypeAccess,
	)
}

func (j *JWTProvider) GenerateRefreshToken(id uuid.UUID) (domain.Token, error) {
	return j.generate(
		id,
		j.config.RefreshTokenTTL,
		domain.TokenTypeRefresh,
	)
}

func (j *JWTProvider) ParseToken(token string) (domain.Claims, error) {
	var claims jwtClaims

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (any, error) {
			return []byte(j.config.Secret), nil
		},
	)
	if err != nil {
		return domain.Claims{}, fmt.Errorf("parse JWT: %w", err)
	}

	domainClaims := domain.NewClaims(claims.UserID, claims.Type)

	return domainClaims, nil
}
