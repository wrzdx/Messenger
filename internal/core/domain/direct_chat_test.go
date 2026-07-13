package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewDirectChatNormalizesUsers(t *testing.T) {
	now := time.Now()
	chatID := uuid.New()
	lowerUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	higherUserID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")

	chat, err := NewDirectChat(chatID, higherUserID, lowerUserID, now)

	require.NoError(t, err)
	require.Equal(t, DirectChat{
		Chat: Chat{
			ID:             chatID,
			LastActivityAt: now,
			CreatedAt:      now,
		},
		User1ID: lowerUserID,
		User2ID: higherUserID,
	}, chat)
}

func TestNewDirectChatReturnsZeroValueWhenInvalid(t *testing.T) {
	now := time.Now()
	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	tests := []struct {
		name      string
		chatID    uuid.UUID
		user1ID   uuid.UUID
		user2ID   uuid.UUID
		createdAt time.Time
		wantError error
	}{
		{name: "nil chat id", chatID: uuid.Nil, user1ID: user1ID, user2ID: user2ID, createdAt: now, wantError: ErrInvalidChat},
		{name: "zero created_at", chatID: uuid.New(), user1ID: user1ID, user2ID: user2ID, wantError: ErrInvalidChat},
		{name: "nil first user id", chatID: uuid.New(), user1ID: uuid.Nil, user2ID: user2ID, createdAt: now, wantError: ErrInvalidDirectChat},
		{name: "nil second user id", chatID: uuid.New(), user1ID: user1ID, user2ID: uuid.Nil, createdAt: now, wantError: ErrInvalidDirectChat},
		{name: "same users", chatID: uuid.New(), user1ID: user1ID, user2ID: user1ID, createdAt: now, wantError: ErrInvalidDirectChat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat, err := NewDirectChat(tt.chatID, tt.user1ID, tt.user2ID, tt.createdAt)

			require.ErrorIs(t, err, tt.wantError)
			require.Zero(t, chat)
		})
	}
}

func TestDirectChatValidate(t *testing.T) {
	now := time.Now()
	lowerUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	higherUserID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	validChat := Chat{ID: uuid.New(), LastActivityAt: now, CreatedAt: now}

	tests := []struct {
		name      string
		chat      DirectChat
		wantError error
	}{
		{
			name: "valid normalized chat",
			chat: DirectChat{Chat: validChat, User1ID: lowerUserID, User2ID: higherUserID},
		},
		{
			name:      "users are not normalized",
			chat:      DirectChat{Chat: validChat, User1ID: higherUserID, User2ID: lowerUserID},
			wantError: ErrInvalidDirectChat,
		},
		{
			name:      "invalid common chat",
			chat:      DirectChat{Chat: Chat{}, User1ID: lowerUserID, User2ID: higherUserID},
			wantError: ErrInvalidChat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.chat.Validate(), tt.wantError)
		})
	}
}
