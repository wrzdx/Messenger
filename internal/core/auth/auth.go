package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrPasswordMismatch     = errors.New("passwords do not match")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrInvalidTokenLifetime = errors.New("invalid token lifetime")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)

type TokenPair struct {
	Access  string
	Refresh string
}

type RefreshTokenClaims struct {
	SessionID uuid.UUID `json:"session_id"`
	TokenID   uuid.UUID `json:"token_id"`
}

type AccessTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
}

func (c RefreshTokenClaims) Validate() error {
	if c.SessionID == uuid.Nil ||
		c.TokenID == uuid.Nil {
		return ErrInvalidClaims
	}
	return nil
}

func (c AccessTokenClaims) Validate() error {
	if c.UserID == uuid.Nil {
		return ErrInvalidClaims
	}
	return nil
}

type TokenLifetime struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (l TokenLifetime) Validate() error {
	if l.ExpiresAt.IsZero() || l.IssuedAt.IsZero() || !l.ExpiresAt.After(l.IssuedAt) {
		return ErrInvalidTokenLifetime
	}
	return nil
}
