package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleMember UserRole = "member"
	UserRoleAdmin  UserRole = "admin"
	UserRoleOwner  UserRole = "owner"
)

var (
	ErrInvalidParticipants = errors.New("direct chat requires exactly one participant")
	ErrInvalidUserRole     = errors.New("user role can be either 'member' or 'admin' or 'owner'")
)

type ChatParticipant struct {
	ChatID            uuid.UUID
	UserID            uuid.UUID
	Role              UserRole
	LastReadMessageID *uuid.UUID
	JoinedAt          time.Time
}

func ValidateUserRole(role string) error {
	switch UserRole(role) {
	case UserRoleMember,
		UserRoleAdmin,
		UserRoleOwner:
		return nil
	default:
		return ErrInvalidUserRole
	}
}
