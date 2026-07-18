package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrInvalidMessage = errors.New("invalid message")

type Message struct {
	ID              uuid.UUID
	ClientMessageID uuid.UUID
	ChatID          uuid.UUID
	SenderID        uuid.UUID
	Content         string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

func NewMessage(
	id, clientMessageID, chatID, senderID uuid.UUID,
	content string,
	createdAt time.Time,
) (Message, error) {
	message := Message{
		ID:              id,
		ClientMessageID: clientMessageID,
		ChatID:          chatID,
		SenderID:        senderID,
		Content:         content,
		CreatedAt:       createdAt,
	}

	message = message.normalize()
	if err := message.Validate(); err != nil {
		return Message{}, err
	}
	return message, nil
}

func (m Message) normalize() Message {
	m.Content = strings.TrimSpace(m.Content)
	return m
}

func (m Message) Validate() error {
	fields := make(map[string]string)

	if m.ID == uuid.Nil {
		return fmt.Errorf("id is nil: %w", ErrInvalidMessage)
	}

	if m.ClientMessageID == uuid.Nil {
		return fmt.Errorf("client_message_id is nil: %w", ErrInvalidMessage)
	}

	if m.ChatID == uuid.Nil {
		return fmt.Errorf("chat_id is nil: %w", ErrInvalidMessage)
	}
	if m.SenderID == uuid.Nil {
		return fmt.Errorf("sender_id is nil: %w", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.Content) != m.Content {
		return fmt.Errorf("content is not normalized: %w", ErrInvalidMessage)
	}
	if m.CreatedAt.IsZero() {
		return fmt.Errorf("created_at is zero value: %w", ErrInvalidMessage)
	}
	if m.UpdatedAt != nil && m.UpdatedAt.IsZero() {
		return fmt.Errorf("updated_at is zero value: %w", ErrInvalidMessage)
	}
	if m.UpdatedAt != nil && m.UpdatedAt.Before(m.CreatedAt) {
		return fmt.Errorf("updated_at before created_at: %w", ErrInvalidMessage)
	}
	if l := utf8.RuneCountInString(m.Content); l < 1 || l > 4096 {
		fields["content"] = "content must contain between 1 and 4096 characters"
	}

	if len(fields) > 0 {
		return DetailedError{
			Err:     ErrInvalidMessage,
			Details: fields,
		}
	}
	return nil
}
