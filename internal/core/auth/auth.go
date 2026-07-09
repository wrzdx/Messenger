package auth

import "errors"

var ErrInvalidToken = errors.New("invalid token")
var ErrPasswordMismatch = errors.New("passwords do not match")

type TokenPair struct {
	Access  string
	Refresh string
}
