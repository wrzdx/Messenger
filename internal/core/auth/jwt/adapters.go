package core_auth_jwt

import (
	"messenger/internal/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtClaims struct {
	UserID    uuid.UUID        `json:"user_id"`
	Type      domain.TokenType `json:"type"`
	IssuedAt  time.Time        `json:"issued_at"`
	ExpiresAt time.Time        `json:"expires_at"`
	jwt.RegisteredClaims
}
