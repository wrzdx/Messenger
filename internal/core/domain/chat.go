package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	MinChatTitleLength int = 1
	MaxChatTitleLength int = 128
)

var (
	ErrInvalidChatType = errors.New("chat type can be either 'direct' or 'group'")
)

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

type Chat struct {
	ID             uuid.UUID
	Type           ChatType
	Title          *string // direct has no name
	LastMessageID  *uuid.UUID
	LastActivityAt time.Time
	CreatedAt      time.Time
}

func NewChat(
	id uuid.UUID,
	chatType ChatType,
	title *string,
	lastMessageID *uuid.UUID,
	lastActivityAt time.Time,
	createdAt time.Time,
) Chat {
	return Chat{
		ID:             id,
		Type:           chatType,
		Title:          title,
		LastMessageID:  lastMessageID,
		LastActivityAt: lastActivityAt,
		CreatedAt:      createdAt,
	}
}

func ValidateChatTitle(chatTitle string) error {
	return validateLength(
		"title",
		chatTitle,
		new(MinChatTitleLength),
		new(MaxChatTitleLength),
	)
}

func ValidateChatType(chatType string) error {
	switch ChatType(chatType) {
	case ChatTypeDirect,
		ChatTypeGroup:
		return nil
	default:
		return ErrInvalidChatType
	}
}
