package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewMessageNormalizesContent(t *testing.T) {
	id := uuid.New()
	clientMessageID := uuid.New()
	chatID := uuid.New()
	senderID := uuid.New()
	createdAt := time.Now()

	message, err := NewMessage(
		id,
		clientMessageID,
		chatID,
		senderID,
		" \n Hello, world! \t",
		createdAt,
	)

	require.NoError(t, err)
	require.Equal(t, Message{
		ID:              id,
		ClientMessageID: clientMessageID,
		ChatID:          chatID,
		SenderID:        senderID,
		Content:         "Hello, world!",
		CreatedAt:       createdAt,
	}, message)
}

func TestNewMessageReturnsZeroValueWhenInvalid(t *testing.T) {
	now := time.Now()
	validID := uuid.New()
	validClientMessageID := uuid.New()
	validChatID := uuid.New()
	validSenderID := uuid.New()

	tests := []struct {
		name            string
		id              uuid.UUID
		clientMessageID uuid.UUID
		chatID          uuid.UUID
		senderID        uuid.UUID
		content         string
		createdAt       time.Time
	}{
		{
			name:            "nil message id",
			clientMessageID: validClientMessageID,
			chatID:          validChatID,
			senderID:        validSenderID,
			content:         "message",
			createdAt:       now,
		},
		{
			name:            "nil client message id",
			id:              validID,
			clientMessageID: uuid.Nil,
			chatID:          validChatID,
			senderID:        validSenderID,
			content:         "message",
			createdAt:       now,
		},
		{
			name:            "nil chat id",
			id:              validID,
			clientMessageID: validClientMessageID,
			senderID:        validSenderID,
			content:         "message",
			createdAt:       now,
		},
		{
			name:            "nil sender id",
			id:              validID,
			clientMessageID: validClientMessageID,
			chatID:          validChatID,
			content:         "message",
			createdAt:       now,
		},
		{
			name:            "empty content after normalization",
			id:              validID,
			clientMessageID: validClientMessageID,
			chatID:          validChatID,
			senderID:        validSenderID,
			content:         " \n\t ",
			createdAt:       now,
		},
		{
			name:            "content longer than 4096 runes",
			id:              validID,
			clientMessageID: validClientMessageID,
			chatID:          validChatID,
			senderID:        validSenderID,
			content:         strings.Repeat("я", 4097),
			createdAt:       now,
		},
		{
			name:            "zero created at",
			id:              validID,
			clientMessageID: validClientMessageID,
			chatID:          validChatID,
			senderID:        validSenderID,
			content:         "message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message, err := NewMessage(
				test.id,
				test.clientMessageID,
				test.chatID,
				test.senderID,
				test.content,
				test.createdAt,
			)

			require.ErrorIs(t, err, ErrInvalidMessage)
			require.Zero(t, message)
		})
	}
}

func TestMessageValidate(t *testing.T) {
	now := time.Now()
	validMessage := func() Message {
		return Message{
			ID:              uuid.New(),
			ClientMessageID: uuid.New(),
			ChatID:          uuid.New(),
			SenderID:        uuid.New(),
			Content:         "message",
			CreatedAt:       now,
		}
	}

	t.Run("accepts updated message", func(t *testing.T) {
		message := validMessage()
		updatedAt := now.Add(time.Second)
		message.UpdatedAt = &updatedAt

		require.NoError(t, message.Validate())
	})

	t.Run("rejects non-normalized content", func(t *testing.T) {
		message := validMessage()
		message.Content = " message "

		require.ErrorIs(t, message.Validate(), ErrInvalidMessage)
	})

	t.Run("rejects zero updated at", func(t *testing.T) {
		message := validMessage()
		message.UpdatedAt = new(time.Time)

		require.ErrorIs(t, message.Validate(), ErrInvalidMessage)
	})

	t.Run("rejects updated at before creation", func(t *testing.T) {
		message := validMessage()
		updatedAt := now.Add(-time.Second)
		message.UpdatedAt = &updatedAt

		require.ErrorIs(t, message.Validate(), ErrInvalidMessage)
	})
}

func TestMessageUpdate(t *testing.T) {
	createdAt := time.Date(2026, time.August, 20, 12, 0, 0, 0, time.UTC)
	message, err := NewMessage(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"old content",
		createdAt,
	)
	require.NoError(t, err)

	t.Run("normalizes content and records update time", func(t *testing.T) {
		updatedAt := createdAt.Add(time.Minute)

		updated, err := message.Update(" \n new content \t", updatedAt)

		require.NoError(t, err)
		require.Equal(t, "new content", updated.Content)
		require.NotNil(t, updated.UpdatedAt)
		require.True(t, updatedAt.Equal(*updated.UpdatedAt))
		require.Equal(t, "old content", message.Content)
		require.Nil(t, message.UpdatedAt)
	})

	for _, testCase := range []struct {
		name      string
		content   string
		updatedAt time.Time
	}{
		{
			name:      "rejects empty normalized content",
			content:   " \n\t ",
			updatedAt: createdAt.Add(time.Minute),
		},
		{
			name:      "rejects content longer than limit",
			content:   strings.Repeat("я", 4097),
			updatedAt: createdAt.Add(time.Minute),
		},
		{
			name:      "rejects update before creation",
			content:   "new content",
			updatedAt: createdAt.Add(-time.Second),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			updated, err := message.Update(testCase.content, testCase.updatedAt)

			require.ErrorIs(t, err, ErrInvalidMessage)
			require.Zero(t, updated)
		})
	}
}
