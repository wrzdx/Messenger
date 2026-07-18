package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrInvalidGroupChat = errors.New("invalid group chat")

type GroupChat struct {
	Chat  Chat
	Title string
}

func NewGroupChat(id uuid.UUID, title string, createdAt time.Time) (GroupChat, error) {
	chat, err := newChat(id, ChatTypeGroup, createdAt)
	if err != nil {
		return GroupChat{}, err
	}
	groupChat := GroupChat{
		Chat:  chat,
		Title: title,
	}
	groupChat = groupChat.normalize()
	if err := groupChat.Validate(); err != nil {
		return GroupChat{}, err
	}
	return groupChat, nil
}

func (c GroupChat) normalize() GroupChat {
	c.Title = strings.TrimSpace(c.Title)
	return c
}

func (c GroupChat) Validate() error {
	fields := make(map[string]string)
	if err := c.Chat.Validate(); err != nil {
		return err
	}
	if c.Chat.Type != ChatTypeGroup {
		return fmt.Errorf("wrong chat type for group chat: %w", ErrInvalidGroupChat)
	}
	if strings.TrimSpace(c.Title) != c.Title {
		return fmt.Errorf("title is not normalized: %w", ErrInvalidGroupChat)
	}
	if l := utf8.RuneCountInString(c.Title); l < 1 || l > 128 {
		fields["title"] += "title must contain between 1 and 128 characters"
	}
	if len(fields) > 0 {
		return DetailedError{
			Err:     ErrInvalidGroupChat,
			Details: fields,
		}
	}
	return nil
}
