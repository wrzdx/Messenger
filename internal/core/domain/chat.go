package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidChat = errors.New("invalid chat")
)

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

type Chat struct {
	ID             uuid.UUID
	Type           ChatType
	LastMessageID  *uuid.UUID
	LastActivityAt time.Time
	CreatedAt      time.Time
}

func newChat(
	id uuid.UUID,
	chatType ChatType,
	createdAt time.Time,
) (Chat, error) {
	chat := Chat{
		ID:             id,
		Type:           chatType,
		LastMessageID:  nil,
		LastActivityAt: createdAt,
		CreatedAt:      createdAt,
	}

	if err := chat.validate(); err != nil {
		return Chat{}, err
	}

	return chat, nil
}

func (c Chat) validate() error {
	if c.ID == uuid.Nil {
		return fmt.Errorf("id is nil: %w", ErrInvalidChat)
	}
	if c.CreatedAt.IsZero() {
		return fmt.Errorf("created_at is zero value: %w", ErrInvalidChat)
	}
	if c.LastActivityAt.IsZero() {
		return fmt.Errorf("last_activity_at is zero value: %w", ErrInvalidChat)
	}
	if c.LastActivityAt.Before(c.CreatedAt) {
		return fmt.Errorf("last_activity_at cannot be before created_at: %w", ErrInvalidChat)
	}
	if c.LastMessageID != nil && *c.LastMessageID == uuid.Nil {
		return fmt.Errorf("last_message_id is nil: %w", ErrInvalidChat)
	}
	if c.Type != ChatTypeDirect && c.Type != ChatTypeGroup {
		return fmt.Errorf("unknown chat type: %w", ErrInvalidChat)
	}
	return nil
}
