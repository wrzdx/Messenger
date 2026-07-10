package domain

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	ErrEmptyMessage = errors.New("missing message content")
	ErrLongMessage  = errors.New("message content can be at most 4096 characters")
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	CreatedAt string
	UpdatedAt *string
}

func ValidateMessageContent(content string) error {
	l := utf8.RuneCountInString(content)
	if l < 1 {
		return ErrEmptyMessage
	}
	if l > 4096 {
		return ErrLongMessage
	}

	return nil
}
