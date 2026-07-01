package core_auth

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
	Compare(hash, password string) error
}
