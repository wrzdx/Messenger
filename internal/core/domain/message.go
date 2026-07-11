package domain

import (
	"github.com/google/uuid"
)

const (
	MinMessageLength int = 1
	MaxMessageLength int = 4096
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
	return validateLength(
		"content",
		content,
		new(MinMessageLength),
		new(MaxMessageLength),
	)
}
