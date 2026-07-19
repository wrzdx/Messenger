package messages_service

import (
	"errors"
	"testing"
	"time"

	"messenger/internal/core/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetMessages(t *testing.T) {
	requesterID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	chatID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	t.Run("uses default limit and returns non-nil empty page", func(t *testing.T) {
		repository := NewMockMessagesRepository(t)
		repository.EXPECT().CheckParticipant(t.Context(), chatID, requesterID).Return(nil)
		repository.EXPECT().GetMessages(t.Context(), chatID, (*MessageCursor)(nil), 51).Return(nil, nil)
		service := newGetMessagesTestService(t, repository)

		page, err := service.GetMessages(t.Context(), requesterID, GetMessagesQuery{ChatID: chatID})

		require.NoError(t, err)
		require.NotNil(t, page.Messages)
		require.Empty(t, page.Messages)
		require.Nil(t, page.NextCursor)
	})

	for _, testCase := range []struct {
		name       string
		resultSize int
		hasMore    bool
	}{
		{name: "returns fewer than limit", resultSize: 1},
		{name: "returns exactly limit", resultSize: 2},
		{name: "removes lookahead message and returns cursor", resultSize: 3, hasMore: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			messages := newGetMessagesTestMessages(t, chatID, requesterID, testCase.resultSize)
			repository := NewMockMessagesRepository(t)
			repository.EXPECT().CheckParticipant(t.Context(), chatID, requesterID).Return(nil)
			repository.EXPECT().GetMessages(t.Context(), chatID, (*MessageCursor)(nil), 3).
				Return(messages, nil)
			service := newGetMessagesTestService(t, repository)

			page, err := service.GetMessages(t.Context(), requesterID, GetMessagesQuery{
				ChatID: chatID,
				Limit:  2,
			})

			require.NoError(t, err)
			expectedSize := testCase.resultSize
			if testCase.hasMore {
				expectedSize = 2
			}
			require.Equal(t, messages[:expectedSize], page.Messages)
			if !testCase.hasMore {
				require.Nil(t, page.NextCursor)
				return
			}
			require.Equal(t, &MessageCursor{
				MessageID: messages[1].ID,
				CreatedAt: messages[1].CreatedAt,
			}, page.NextCursor)
			require.NotEqual(t, messages[2].ID, page.NextCursor.MessageID)
		})
	}

	t.Run("forwards cursor to repository", func(t *testing.T) {
		before := &MessageCursor{
			MessageID: uuid.New(),
			CreatedAt: time.Date(2026, time.July, 19, 12, 0, 0, 0, time.UTC),
		}
		repository := NewMockMessagesRepository(t)
		repository.EXPECT().CheckParticipant(t.Context(), chatID, requesterID).Return(nil)
		repository.EXPECT().GetMessages(t.Context(), chatID, before, 11).
			Return([]domain.Message{}, nil)
		service := newGetMessagesTestService(t, repository)

		page, err := service.GetMessages(t.Context(), requesterID, GetMessagesQuery{
			ChatID: chatID,
			Before: before,
			Limit:  10,
		})

		require.NoError(t, err)
		require.Empty(t, page.Messages)
		require.Nil(t, page.NextCursor)
	})

	t.Run("returns participant lookup error without loading messages", func(t *testing.T) {
		lookupErr := errors.New("participant lookup failed")
		repository := NewMockMessagesRepository(t)
		repository.EXPECT().CheckParticipant(t.Context(), chatID, requesterID).Return(lookupErr)
		service := newGetMessagesTestService(t, repository)

		page, err := service.GetMessages(t.Context(), requesterID, GetMessagesQuery{
			ChatID: chatID,
			Limit:  10,
		})

		require.ErrorIs(t, err, lookupErr)
		require.Zero(t, page)
	})

	t.Run("returns message repository error", func(t *testing.T) {
		loadErr := errors.New("message lookup failed")
		repository := NewMockMessagesRepository(t)
		repository.EXPECT().CheckParticipant(t.Context(), chatID, requesterID).Return(nil)
		repository.EXPECT().GetMessages(t.Context(), chatID, (*MessageCursor)(nil), 11).
			Return(nil, loadErr)
		service := newGetMessagesTestService(t, repository)

		page, err := service.GetMessages(t.Context(), requesterID, GetMessagesQuery{
			ChatID: chatID,
			Limit:  10,
		})

		require.ErrorIs(t, err, loadErr)
		require.Zero(t, page)
	})

	t.Run("hides nil requester as not found", func(t *testing.T) {
		service := newGetMessagesTestService(t, NewMockMessagesRepository(t))

		page, err := service.GetMessages(t.Context(), uuid.Nil, GetMessagesQuery{
			ChatID: chatID,
			Limit:  10,
		})

		require.ErrorIs(t, err, domain.ErrNotFound)
		require.Zero(t, page)
	})

	for _, testCase := range []struct {
		name           string
		query          GetMessagesQuery
		expectedFields map[string]string
	}{
		{
			name:  "rejects nil chat id",
			query: GetMessagesQuery{Limit: 10},
			expectedFields: map[string]string{
				"chat_id": "chat_id is nil",
			},
		},
		{
			name:  "rejects negative limit",
			query: GetMessagesQuery{ChatID: chatID, Limit: -1},
			expectedFields: map[string]string{
				"limit": "limit must be between 1 and 100",
			},
		},
		{
			name:  "rejects limit above maximum",
			query: GetMessagesQuery{ChatID: chatID, Limit: 101},
			expectedFields: map[string]string{
				"limit": "limit must be between 1 and 100",
			},
		},
		{
			name: "rejects invalid cursor",
			query: GetMessagesQuery{
				ChatID: chatID,
				Limit:  10,
				Before: &MessageCursor{},
			},
			expectedFields: map[string]string{
				"created_at": "created_at of message cursor cannot be zero value",
				"message_id": "message id of message cursor cannot be nil",
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			service := newGetMessagesTestService(t, NewMockMessagesRepository(t))

			page, err := service.GetMessages(t.Context(), requesterID, testCase.query)

			require.ErrorIs(t, err, ErrInvalidGetMessagesQuery)
			require.Zero(t, page)
			var detailed domain.DetailedError
			require.ErrorAs(t, err, &detailed)
			require.Equal(t, testCase.expectedFields, detailed.Fields())
		})
	}
}

func newGetMessagesTestService(
	t *testing.T,
	repository MessagesRepository,
) *MessagesService {
	t.Helper()
	return NewMessagesService(
		repository,
		NewMockChatsRepository(t),
		NewMockTXManager(t),
	)
}

func newGetMessagesTestMessages(
	t *testing.T,
	chatID, senderID uuid.UUID,
	count int,
) []domain.Message {
	t.Helper()
	messages := make([]domain.Message, 0, count)
	createdAt := time.Date(2026, time.July, 19, 12, 0, 0, 0, time.UTC)
	for i := range count {
		message, err := domain.NewMessage(
			uuid.New(),
			uuid.New(),
			chatID,
			senderID,
			"Message",
			createdAt.Add(-time.Duration(i)*time.Minute),
		)
		require.NoError(t, err)
		messages = append(messages, message)
	}
	return messages
}
