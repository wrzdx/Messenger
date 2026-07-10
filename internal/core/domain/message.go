package domain

import "github.com/google/uuid"

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	CreatedAt string
	UpdatedAt *string
}
