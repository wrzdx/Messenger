package domain

import (
	"time"

	"github.com/google/uuid"
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
