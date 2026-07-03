package core_auth

import (
	"time"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type JWTProvider interface {
	GenerateAccessToken(id int) (string, time.Time, error)
	GenerateRefreshToken(id int) (string, time.Time, error)
	ParseToken(token string) (Claims, error)
}

type Claims struct {
	UserID    int       `json:"user_id"`
	Type      TokenType `json:"type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
