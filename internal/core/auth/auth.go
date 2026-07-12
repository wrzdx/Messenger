package auth

import "errors"

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrPasswordMismatch = errors.New("passwords do not match")

type TokenPair struct {
	Access  string
	Refresh string
}
