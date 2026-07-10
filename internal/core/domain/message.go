package domain

import (
	"errors"

	"github.com/google/uuid"
)


var (
	ErrEmptyMessage = errors.New("missing message content")
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	CreatedAt string
	UpdatedAt *string
}
