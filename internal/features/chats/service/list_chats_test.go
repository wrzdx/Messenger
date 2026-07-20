package chats_service

import (
	"errors"
	"messenger/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListChats(t *testing.T) {
	requesterID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	peerID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	baseTime := time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)

	t.Run("returns empty page with default limit", func(t *testing.T) {
		repository := NewMockChatsRepository(t)
		repository.EXPECT().
			ListChats(mock.Anything, requesterID, (*ChatCursor)(nil), 51).
			Return(nil, nil)
		service := NewChatsService(repository, nil, nil)

		page, err := service.ListChats(t.Context(), requesterID, ListChatsQuery{})

		require.NoError(t, err)
		require.NotNil(t, page.Chats)
		require.Empty(t, page.Chats)
		require.Nil(t, page.NextCursor)
	})

	t.Run("returns page and cursor from last visible chat", func(t *testing.T) {
		first := newListChatsDirectItem(t, requesterID, peerID, baseTime.Add(2*time.Minute), true)
		second := newListChatsGroupItem(t, requesterID, baseTime.Add(time.Minute), false)
		extra := newListChatsGroupItem(t, requesterID, baseTime, false)
		before := &ChatCursor{ChatID: uuid.New(), LastActivityAt: baseTime.Add(time.Hour)}
		repository := NewMockChatsRepository(t)
		repository.EXPECT().
			ListChats(mock.Anything, requesterID, before, 3).
			Return([]ChatItem{first, second, extra}, nil)
		service := NewChatsService(repository, nil, nil)

		page, err := service.ListChats(t.Context(), requesterID, ListChatsQuery{
			Before: before,
			Limit:  2,
		})

		require.NoError(t, err)
		require.Equal(t, []ChatItem{first, second}, page.Chats)
		require.NotNil(t, page.NextCursor)
		require.Equal(t, second.Group.Chat.ID, page.NextCursor.ChatID)
		require.True(t, second.Group.Chat.LastActivityAt.Equal(page.NextCursor.LastActivityAt))
	})

	t.Run("rejects nil requester", func(t *testing.T) {
		service := NewChatsService(NewMockChatsRepository(t), nil, nil)

		page, err := service.ListChats(t.Context(), uuid.Nil, ListChatsQuery{})

		require.ErrorIs(t, err, domain.ErrNotFound)
		require.Empty(t, page)
	})

	testCases := []struct {
		name  string
		query ListChatsQuery
	}{
		{name: "negative limit", query: ListChatsQuery{Limit: -1}},
		{name: "limit above maximum", query: ListChatsQuery{Limit: 101}},
		{
			name: "nil cursor chat id",
			query: ListChatsQuery{Before: &ChatCursor{
				LastActivityAt: baseTime,
			}},
		},
		{
			name: "zero cursor activity",
			query: ListChatsQuery{Before: &ChatCursor{
				ChatID: uuid.New(),
			}},
		},
	}
	for _, testCase := range testCases {
		t.Run("rejects "+testCase.name, func(t *testing.T) {
			service := NewChatsService(NewMockChatsRepository(t), nil, nil)

			page, err := service.ListChats(t.Context(), requesterID, testCase.query)

			require.ErrorIs(t, err, ErrInvalidListChatsQuery)
			require.Empty(t, page)
		})
	}

	t.Run("returns repository error", func(t *testing.T) {
		repositoryErr := errors.New("database unavailable")
		repository := NewMockChatsRepository(t)
		repository.EXPECT().
			ListChats(mock.Anything, requesterID, (*ChatCursor)(nil), 51).
			Return(nil, repositoryErr)
		service := NewChatsService(repository, nil, nil)

		page, err := service.ListChats(t.Context(), requesterID, ListChatsQuery{})

		require.ErrorIs(t, err, repositoryErr)
		require.Empty(t, page)
	})
}

func TestChatItemValidate(t *testing.T) {
	requesterID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	peerID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	createdAt := time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)

	t.Run("accepts valid direct with last message", func(t *testing.T) {
		item := newListChatsDirectItem(t, requesterID, peerID, createdAt, true)

		require.NoError(t, item.Validate())
	})

	t.Run("accepts valid empty group", func(t *testing.T) {
		item := newListChatsGroupItem(t, requesterID, createdAt, false)

		require.NoError(t, item.Validate())
	})

	t.Run("rejects missing subtype", func(t *testing.T) {
		require.ErrorIs(t, (ChatItem{}).Validate(), ErrInvalidChatItem)
	})

	t.Run("rejects both subtypes", func(t *testing.T) {
		direct := newListChatsDirectItem(t, requesterID, peerID, createdAt, false)
		group := newListChatsGroupItem(t, requesterID, createdAt, false)
		direct.Group = group.Group

		require.ErrorIs(t, direct.Validate(), ErrInvalidChatItem)
	})

	t.Run("rejects mismatched last message", func(t *testing.T) {
		item := newListChatsDirectItem(t, requesterID, peerID, createdAt, true)
		item.LastMessage.Message.ChatID = uuid.New()

		require.ErrorIs(t, item.Validate(), ErrInvalidChatItem)
	})
}

func newListChatsDirectItem(
	t *testing.T,
	requesterID uuid.UUID,
	peerID uuid.UUID,
	activityAt time.Time,
	withMessage bool,
) ChatItem {
	t.Helper()
	direct, err := domain.NewDirectChat(uuid.New(), requesterID, peerID, activityAt.Add(-time.Hour))
	require.NoError(t, err)
	peerProfile, err := domain.NewUserProfile("Peer_123", "Peer", nil, nil)
	require.NoError(t, err)
	item := ChatItem{Direct: &DirectChatItem{
		Chat:        direct,
		PeerID:      peerID,
		PeerProfile: peerProfile,
	}}
	if withMessage {
		message, err := domain.NewMessage(
			uuid.New(),
			uuid.New(),
			direct.Chat.ID,
			requesterID,
			"hello",
			activityAt,
		)
		require.NoError(t, err)
		item.Direct.Chat.Chat.LastMessageID = &message.ID
		item.Direct.Chat.Chat.LastActivityAt = message.CreatedAt
		senderProfile, err := domain.NewUserProfile("Sender_123", "Sender", nil, nil)
		require.NoError(t, err)
		item.LastMessage = &LastMessageItem{Message: message, SenderProfile: senderProfile}
	} else {
		item.Direct.Chat.Chat.LastActivityAt = activityAt
	}
	require.NoError(t, item.Validate())
	return item
}

func newListChatsGroupItem(
	t *testing.T,
	senderID uuid.UUID,
	activityAt time.Time,
	withMessage bool,
) ChatItem {
	t.Helper()
	group, err := domain.NewGroupChat(uuid.New(), "Backend", activityAt.Add(-time.Hour))
	require.NoError(t, err)
	item := ChatItem{Group: &group}
	if withMessage {
		message, err := domain.NewMessage(
			uuid.New(),
			uuid.New(),
			group.Chat.ID,
			senderID,
			"group message",
			activityAt,
		)
		require.NoError(t, err)
		item.Group.Chat.LastMessageID = &message.ID
		item.Group.Chat.LastActivityAt = message.CreatedAt
		senderProfile, err := domain.NewUserProfile("Author_123", "Author", nil, nil)
		require.NoError(t, err)
		item.LastMessage = &LastMessageItem{Message: message, SenderProfile: senderProfile}
	} else {
		item.Group.Chat.LastActivityAt = activityAt
	}
	require.NoError(t, item.Validate())
	return item
}
