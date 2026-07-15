package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewChat(t *testing.T) {
	now := time.Now()
	id := uuid.New()

	chat, err := newChat(id, now)

	require.NoError(t, err)
	require.Equal(t, Chat{
		ID:             id,
		LastActivityAt: now,
		CreatedAt:      now,
	}, chat)
}

func TestNewChatReturnsZeroValueWhenInvalid(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		createdAt time.Time
	}{
		{name: "nil id", id: uuid.Nil, createdAt: time.Now()},
		{name: "zero created_at", id: uuid.New(), createdAt: time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat, err := newChat(tt.id, tt.createdAt)

			require.ErrorIs(t, err, ErrInvalidChat)
			require.Zero(t, chat)
		})
	}
}

func TestChatValidate(t *testing.T) {
	now := time.Now()
	lastMessageID := uuid.New()

	validChat := func() Chat {
		return Chat{
			ID:             uuid.New(),
			LastActivityAt: now,
			CreatedAt:      now,
		}
	}

	tests := []struct {
		name      string
		change    func(*Chat)
		wantError error
	}{
		{name: "valid initial chat", change: func(*Chat) {}},
		{
			name: "valid active chat",
			change: func(chat *Chat) {
				chat.LastMessageID = &lastMessageID
				chat.LastActivityAt = now.Add(time.Hour)
			},
		},
		{
			name: "nil id",
			change: func(chat *Chat) {
				chat.ID = uuid.Nil
			},
			wantError: ErrInvalidChat,
		},
		{
			name: "zero created_at",
			change: func(chat *Chat) {
				chat.CreatedAt = time.Time{}
			},
			wantError: ErrInvalidChat,
		},
		{
			name: "zero last_activity_at",
			change: func(chat *Chat) {
				chat.LastActivityAt = time.Time{}
			},
			wantError: ErrInvalidChat,
		},
		{
			name: "last_activity_at before created_at",
			change: func(chat *Chat) {
				chat.LastActivityAt = now.Add(-time.Second)
			},
			wantError: ErrInvalidChat,
		},
		{
			name: "nil last_message_id value",
			change: func(chat *Chat) {
				chat.LastMessageID = new(uuid.UUID)
			},
			wantError: ErrInvalidChat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := validChat()
			tt.change(&chat)

			require.ErrorIs(t, chat.validate(), tt.wantError)
		})
	}
}
