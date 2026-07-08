package auth

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type AccessClaims struct {
	UserID uuid.UUID
}

func NewAccessClaims(userID uuid.UUID) AccessClaims {
	return AccessClaims{
		UserID: userID,
	}
}

type RefreshClaims struct {
	UserID  uuid.UUID
	TokenID uuid.UUID
}

func NewRefreshClaims(userID, tokenID uuid.UUID) RefreshClaims {
	return RefreshClaims{
		UserID:  userID,
		TokenID: tokenID,
	}
}

type TokenPair struct {
	Access  string
	Refresh string
}

func NewTokenPair(access, refresh string) TokenPair {
	return TokenPair{
		Access:  access,
		Refresh: refresh,
	}
}

type TokenService interface {
	GenerateTokenPair(user AccessClaims, tokenID uuid.UUID) (TokenPair, error)
	ParseAccessToken(token string) (AccessClaims, error)
	ParseRefreshToken(token string) (RefreshClaims, error)
}
