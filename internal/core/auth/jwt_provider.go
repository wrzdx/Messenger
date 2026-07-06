package core_auth

import (
	"time"
)

type claimsKey struct{}
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
