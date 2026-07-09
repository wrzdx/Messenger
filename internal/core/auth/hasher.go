package auth

import "errors"

var ErrPasswordMismatch = errors.New("passwords do not match")

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
