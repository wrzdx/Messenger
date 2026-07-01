package core_auth

import (
	"time"
)

type AccessToken string

type RefreshToken string

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type JWTProvider interface {
	GenerateAccessToken(id int) (AccessToken, time.Time, error)
	GenerateRefreshToken(id int) (RefreshToken, time.Time, error)
	ParseToken(token string) (Claims, error)
}

type Claims struct {
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
