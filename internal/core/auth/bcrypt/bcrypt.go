package auth_bcrypt

import (
	"errors"
	"fmt"
	"messenger/internal/core/auth"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
}

func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{}
}

func (h BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}

	return string(hash), nil
}

func (h BcryptHasher) Compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return auth.ErrPasswordMismatch
	}
	return err
}
