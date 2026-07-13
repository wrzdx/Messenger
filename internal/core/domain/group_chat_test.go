package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewGroupChatNormalizesTitle(t *testing.T) {
	now := time.Now()
	chatID := uuid.New()

	chat, err := NewGroupChat(chatID, "  Group title  ", now)

	require.NoError(t, err)
	require.Equal(t, GroupChat{
		Chat: Chat{
			ID:             chatID,
			LastActivityAt: now,
			CreatedAt:      now,
		},
		Title: "Group title",
	}, chat)
}

func TestNewGroupChatTitleBoundaries(t *testing.T) {
	now := time.Now()

	validTitles := []string{"я", strings.Repeat("я", 128)}
	for _, title := range validTitles {
		chat, err := NewGroupChat(uuid.New(), title, now)
		require.NoError(t, err)
		require.Equal(t, title, chat.Title)
	}

	invalidTitles := []string{"", "   ", strings.Repeat("я", 129)}
	for _, title := range invalidTitles {
		chat, err := NewGroupChat(uuid.New(), title, now)
		require.ErrorIs(t, err, ErrInvalidGroupChat)
		require.Zero(t, chat)

		var detailedError DetailedError
		require.ErrorAs(t, err, &detailedError)
		require.Contains(t, detailedError.Details, "title")
	}
}

func TestNewGroupChatReturnsZeroValueWhenCommonChatIsInvalid(t *testing.T) {
	chat, err := NewGroupChat(uuid.Nil, "Group title", time.Now())

	require.ErrorIs(t, err, ErrInvalidChat)
	require.Zero(t, chat)
}

func TestGroupChatValidate(t *testing.T) {
	now := time.Now()
	validChat := Chat{ID: uuid.New(), LastActivityAt: now, CreatedAt: now}

	tests := []struct {
		name      string
		chat      GroupChat
		wantError error
	}{
		{
			name: "valid group chat",
			chat: GroupChat{Chat: validChat, Title: "Group title"},
		},
		{
			name:      "title is not normalized",
			chat:      GroupChat{Chat: validChat, Title: "  Group title  "},
			wantError: ErrInvalidGroupChat,
		},
		{
			name:      "invalid common chat",
			chat:      GroupChat{Chat: Chat{}, Title: "Group title"},
			wantError: ErrInvalidChat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.chat.Validate(), tt.wantError)
		})
	}
}
