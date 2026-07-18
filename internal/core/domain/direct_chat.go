package domain

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidDirectChat = errors.New("invalid direct chat")

type DirectChat struct {
	Chat    Chat
	User1ID uuid.UUID
	User2ID uuid.UUID
}

func NewDirectChat(
	id, user1ID, user2ID uuid.UUID,
	createdAt time.Time,
) (DirectChat, error) {
	chat, err := newChat(id, ChatTypeDirect, createdAt)
	if err != nil {
		return DirectChat{}, err
	}
	directChat := DirectChat{
		Chat:    chat,
		User1ID: user1ID,
		User2ID: user2ID,
	}
	directChat = directChat.normalize()
	if err := directChat.Validate(); err != nil {
		return DirectChat{}, err
	}
	return directChat, nil
}
func (c DirectChat) normalize() DirectChat {
	if bytes.Compare(c.User1ID[:], c.User2ID[:]) > 0 {
		c.User1ID, c.User2ID = c.User2ID, c.User1ID
	}
	return c
}
func (c DirectChat) Validate() error {
	if err := c.Chat.validate(); err != nil {
		return err
	}
	if c.Chat.Type != ChatTypeDirect {
		return fmt.Errorf("wrong chat type for direct: %w", ErrInvalidDirectChat)
	}
	if c.User1ID == uuid.Nil || c.User2ID == uuid.Nil {
		return fmt.Errorf("id is nil: %w", ErrInvalidDirectChat)
	}
	if c.User1ID == c.User2ID {
		return fmt.Errorf("user ids must be different: %w", ErrInvalidDirectChat)
	}
	if bytes.Compare(c.User1ID[:], c.User2ID[:]) > 0 {
		return fmt.Errorf("user ids are not normalized: %w", ErrInvalidDirectChat)
	}
	return nil
}
