package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewChatParticipant(t *testing.T) {
	chatID := uuid.New()
	userID := uuid.New()
	joinedAt := time.Now()

	t.Run("creates participant without read message", func(t *testing.T) {
		participant, err := NewChatParticipant(chatID, userID, nil, joinedAt)

		require.NoError(t, err)
		require.Equal(t, ChatParticipant{
			ChatID:   chatID,
			UserID:   userID,
			JoinedAt: joinedAt,
		}, participant)
	})

	t.Run("creates participant with read message", func(t *testing.T) {
		lastReadMessageID := uuid.New()

		participant, err := NewChatParticipant(chatID, userID, &lastReadMessageID, joinedAt)

		require.NoError(t, err)
		require.Equal(t, &lastReadMessageID, participant.LastReadMessageID)
	})
}

func TestNewChatParticipantReturnsZeroValueWhenInvalid(t *testing.T) {
	now := time.Now()
	validChatID := uuid.New()
	validUserID := uuid.New()
	nilMessageID := uuid.Nil

	tests := []struct {
		name              string
		chatID            uuid.UUID
		userID            uuid.UUID
		lastReadMessageID *uuid.UUID
		joinedAt          time.Time
	}{
		{
			name:     "nil chat id",
			chatID:   uuid.Nil,
			userID:   validUserID,
			joinedAt: now,
		},
		{
			name:     "nil user id",
			chatID:   validChatID,
			userID:   uuid.Nil,
			joinedAt: now,
		},
		{
			name:              "nil last_read_message_id value",
			chatID:            validChatID,
			userID:            validUserID,
			lastReadMessageID: &nilMessageID,
			joinedAt:          now,
		},
		{
			name:     "zero joined_at",
			chatID:   validChatID,
			userID:   validUserID,
			joinedAt: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			participant, err := NewChatParticipant(
				tt.chatID,
				tt.userID,
				tt.lastReadMessageID,
				tt.joinedAt,
			)

			require.ErrorIs(t, err, ErrInvalidChatParticipant)
			require.Zero(t, participant)
		})
	}
}

func TestNewGroupParticipant(t *testing.T) {
	participant, err := NewChatParticipant(uuid.New(), uuid.New(), nil, time.Now())
	require.NoError(t, err)

	roles := []GroupRole{MemberRole, AdminRole, OwnerRole}
	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			groupParticipant, err := NewGroupParticipant(
				participant.ChatID,
				participant.UserID,
				participant.LastReadMessageID,
				participant.JoinedAt,
				role,
			)

			require.NoError(t, err)
			require.Equal(t, participant, groupParticipant.ChatParticipant)
			require.Equal(t, role, groupParticipant.Role())
		})
	}
}

func TestNewGroupParticipantReturnsZeroValueWhenInvalid(t *testing.T) {
	validParticipant, err := NewChatParticipant(uuid.New(), uuid.New(), nil, time.Now())
	require.NoError(t, err)

	t.Run("invalid common participant", func(t *testing.T) {
		groupParticipant, err := NewGroupParticipant(
			uuid.Nil,
			uuid.Nil,
			nil,
			time.Time{},
			MemberRole,
		)

		require.ErrorIs(t, err, ErrInvalidChatParticipant)
		require.NotErrorIs(t, err, ErrInvalidGroupParticipant)
		require.Zero(t, groupParticipant)
	})

	t.Run("unknown group role", func(t *testing.T) {
		groupParticipant, err := NewGroupParticipant(
			validParticipant.ChatID,
			validParticipant.UserID,
			validParticipant.LastReadMessageID,
			validParticipant.JoinedAt,
			GroupRole("unknown"),
		)

		require.ErrorIs(t, err, ErrInvalidGroupParticipant)
		require.NotErrorIs(t, err, ErrInvalidChatParticipant)
		require.Zero(t, groupParticipant)
	})
}
