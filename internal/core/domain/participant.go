package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidChatParticipant  = errors.New("invalid chat participant")
	ErrInvalidGroupParticipant = errors.New("invalid group participant")
)

type GroupRole string

const (
	MemberRole GroupRole = "member"
	AdminRole  GroupRole = "admin"
	OwnerRole  GroupRole = "owner"
)

type ChatParticipant struct {
	ChatID            uuid.UUID
	UserID            uuid.UUID
	LastReadMessageID *uuid.UUID
	JoinedAt          time.Time
}

type GroupParticipant struct {
	ChatParticipant
	role GroupRole
}

func (p GroupParticipant) Role() GroupRole { return p.role }

func NewChatParticipant(
	chatID, userID uuid.UUID,
	lastReadMessageID *uuid.UUID,
	joinedAt time.Time,
) (ChatParticipant, error) {
	participant := ChatParticipant{
		ChatID:            chatID,
		UserID:            userID,
		LastReadMessageID: lastReadMessageID,
		JoinedAt:          joinedAt,
	}

	if err := participant.Validate(); err != nil {
		return ChatParticipant{}, err
	}

	return participant, nil
}

func (p ChatParticipant) Validate() error {
	if p.ChatID == uuid.Nil {
		return fmt.Errorf("chat_id is nil: %w", ErrInvalidChatParticipant)
	}
	if p.UserID == uuid.Nil {
		return fmt.Errorf("user_id is nil: %w", ErrInvalidChatParticipant)
	}
	if p.LastReadMessageID != nil && *p.LastReadMessageID == uuid.Nil {
		return fmt.Errorf("last_read_message_id is nil: %w", ErrInvalidChatParticipant)
	}
	if p.JoinedAt.IsZero() {
		return fmt.Errorf("joined_at is zero value: %w", ErrInvalidChatParticipant)
	}

	return nil
}

func NewGroupParticipant(
	participant ChatParticipant,
	role GroupRole,
) (GroupParticipant, error) {
	groupParticipant := GroupParticipant{
		ChatParticipant: participant,
		role:            role,
	}

	if err := groupParticipant.Validate(); err != nil {
		return GroupParticipant{}, err
	}

	return groupParticipant, nil
}

func (p GroupParticipant) Validate() error {
	if err := p.ChatParticipant.Validate(); err != nil {
		return err
	}
	if p.role != MemberRole && p.role != AdminRole && p.role != OwnerRole {
		return fmt.Errorf("unknown role: %w", ErrInvalidGroupParticipant)
	}

	return nil
}
