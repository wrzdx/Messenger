package domain

import (
	"errors"
	"fmt"
)

type Entity string

const (
	UserEntity    Entity = "user"
	ChatEntity    Entity = "chat"
	MessageEntity Entity = "message"
)

var (
	// Generic

	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

func NotFoundErr(entity Entity, field, value string) error {
	return fmt.Errorf("%s with %s='%s' %w", entity, field, value, ErrNotFound)
}

func AlreadyExistsErr(entity Entity, field, value string) error {
	return fmt.Errorf("%s with %s='%s' %w", entity, field, value, ErrAlreadyExists)
}
