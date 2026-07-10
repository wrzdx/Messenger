package domain

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	ErrInvalidChatName = errors.New("chat name must be between 1 and 128 characters")
	ErrInvalidUserRole = errors.New("user role can be either 'member' or 'admin' or 'owner'")
)

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

type UserRole string

const (
	UserRoleMember UserRole = "member"
	UserRoleAdmin  UserRole = "admin"
	UserRoleOwner  UserRole = "owner"
)

type Chat struct {
	ID             uuid.UUID
	Type           ChatType
	Name           *string // direct has no name
	LastMessageID  *uuid.UUID
	LastActivityAt time.Time
	CreatedAt      time.Time
}

type ChatParticipant struct {
	ChatID            uuid.UUID
	UserID            uuid.UUID
	Role              UserRole
	LastReadMessageID *uuid.UUID
	JoinedAt          time.Time
}

func ValidateChatName(chatName string) error {
	if l := utf8.RuneCountInString(chatName); l < 1 || l > 128 {
		return ErrInvalidChatName
	}
	return nil
}
