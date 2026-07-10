package domain

import "errors"

var (
	ErrNegativeLimit  = errors.New("limit must be non-negative")
	ErrNegativeOffset = errors.New("offset must be non-negative")
)

func ValidateLimit(limit int) error {
	if limit < 0 {
		return ErrNegativeLimit
	}

	return nil
}

func ValidateOffset(offset int) error {
	if offset < 0 {
		return ErrNegativeOffset
	}

	return nil
}
