//go:build integration

package chats_postgres_repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"
	chats_service "messenger/internal/features/chats/service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestListChats(t *testing.T) {
	config := postgres.NewConfigMust()
	pool, err := postgres.NewPool(t.Context(), config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	t.Run("returns mixed chat page with profiles and last message", func(t *testing.T) {
		fixture := newListChatsRepositoryFixture(t, pool, config.Timeout)

		chats, err := fixture.repository.ListChats(
			t.Context(),
			fixture.requesterID,
			nil,
			10,
		)

		require.NoError(t, err)
		require.Len(t, chats, 2)

		groupItem := chats[0]
		require.Nil(t, groupItem.Direct)
		require.NotNil(t, groupItem.Group)
		require.Equal(t, fixture.group.Chat.ID, groupItem.Group.Chat.ID)
		require.Equal(t, fixture.group.Title, groupItem.Group.Title)
		require.NotNil(t, groupItem.LastMessage)
		require.Equal(t, fixture.message.ID, groupItem.LastMessage.Message.ID)
		require.Equal(t, fixture.message.ClientMessageID, groupItem.LastMessage.Message.ClientMessageID)
		require.Equal(t, fixture.message.ChatID, groupItem.LastMessage.Message.ChatID)
		require.Equal(t, fixture.message.SenderID, groupItem.LastMessage.Message.SenderID)
		require.Equal(t, fixture.message.Content, groupItem.LastMessage.Message.Content)
		require.True(t, fixture.message.CreatedAt.Equal(groupItem.LastMessage.Message.CreatedAt))
		require.Equal(t, fixture.message.UpdatedAt, groupItem.LastMessage.Message.UpdatedAt)
		require.Equal(
			t,
			listChatsUsername(fixture.requesterID),
			groupItem.LastMessage.SenderProfile.Username(),
		)

		directItem := chats[1]
		require.NotNil(t, directItem.Direct)
		require.Nil(t, directItem.Group)
		require.Equal(t, fixture.direct.Chat.ID, directItem.Direct.Chat.Chat.ID)
		require.Equal(t, fixture.peerID, directItem.Direct.PeerID)
		require.Equal(t, fixture.deletedPeerUsername, directItem.Direct.PeerProfile.Username())
		require.NotNil(t, directItem.Direct.PeerDeletedAt)
		require.True(t, fixture.peerDeletedAt.Equal(*directItem.Direct.PeerDeletedAt))
		require.Nil(t, directItem.LastMessage)

		firstPage, err := fixture.repository.ListChats(
			t.Context(), fixture.requesterID, nil, 1,
		)
		require.NoError(t, err)
		require.Len(t, firstPage, 1)
		require.Equal(t, fixture.group.Chat.ID, firstPage[0].Group.Chat.ID)

		secondPage, err := fixture.repository.ListChats(
			t.Context(),
			fixture.requesterID,
			&chats_service.ChatCursor{
				ChatID:         fixture.group.Chat.ID,
				LastActivityAt: fixture.group.Chat.LastActivityAt,
			},
			1,
		)
		require.NoError(t, err)
		require.Len(t, secondPage, 1)
		require.Equal(t, fixture.direct.Chat.ID, secondPage[0].Direct.Chat.Chat.ID)
	})

	t.Run("rejects invalid group state restored from database", func(t *testing.T) {
		fixture := newListChatsRepositoryFixture(t, pool, config.Timeout)
		_, err := pool.Exec(
			t.Context(),
			`UPDATE groups SET title = ' Invalid title ' WHERE chat_id = $1`,
			fixture.group.Chat.ID,
		)
		require.NoError(t, err)

		chats, err := fixture.repository.ListChats(
			t.Context(), fixture.requesterID, nil, 10,
		)

		require.ErrorIs(t, err, domain.ErrInvalidGroupChat)
		require.Nil(t, chats)
	})
}

type listChatsRepositoryFixture struct {
	repository          *ChatsRepository
	requesterID         uuid.UUID
	peerID              uuid.UUID
	direct              domain.DirectChat
	group               domain.GroupChat
	message             domain.Message
	peerDeletedAt       time.Time
	deletedPeerUsername string
}

func newListChatsRepositoryFixture(
	t *testing.T,
	pool *pgxpool.Pool,
	timeout time.Duration,
) listChatsRepositoryFixture {
	t.Helper()

	requesterID := uuid.New()
	peerID := uuid.New()
	insertListChatsUser(t, pool, requesterID)
	insertListChatsUser(t, pool, peerID)

	direct, participant1, participant2 := newCreateDirectAggregate(t, requesterID, peerID)
	repository := NewChatsRepository(pool, timeout)
	require.NoError(t, repository.CreateDirect(t.Context(), direct, participant1, participant2))

	groupCreatedAt := direct.Chat.CreatedAt.Add(time.Minute)
	group, err := domain.NewGroupChat(uuid.New(), "Backend", groupCreatedAt)
	require.NoError(t, err)
	insertListChatsGroup(t, pool, group, requesterID, peerID)

	message, err := domain.NewMessage(
		uuid.New(),
		uuid.New(),
		group.Chat.ID,
		requesterID,
		"hello group",
		groupCreatedAt.Add(time.Minute),
	)
	require.NoError(t, err)
	insertListChatsMessage(t, pool, message)
	_, err = pool.Exec(t.Context(), `
		UPDATE chats
		SET last_message_id = $1, last_activity_at = $2
		WHERE id = $3
	`, message.ID, message.CreatedAt, group.Chat.ID)
	require.NoError(t, err)
	group.Chat.LastMessageID = &message.ID
	group.Chat.LastActivityAt = message.CreatedAt

	outsiderGroup, err := domain.NewGroupChat(
		uuid.New(), "Not visible", message.CreatedAt.Add(time.Minute),
	)
	require.NoError(t, err)
	insertListChatsGroup(t, pool, outsiderGroup, peerID)

	peerDeletedAt := direct.Chat.CreatedAt.Add(2 * time.Minute)
	deletedPeerUsername := "deleted_" + strings.ReplaceAll(peerID.String(), "-", "")[:16]
	_, err = pool.Exec(t.Context(), `
		UPDATE users
		SET username = $1, first_name = 'Deleted Account',
		    last_name = NULL, bio = NULL, deleted_at = $2
		WHERE id = $3
	`, deletedPeerUsername, peerDeletedAt, peerID)
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := pool.Exec(
			ctx,
			`DELETE FROM chats WHERE id IN ($1, $2, $3)`,
			direct.Chat.ID,
			group.Chat.ID,
			outsiderGroup.Chat.ID,
		)
		require.NoError(t, err)
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, requesterID, peerID)
		require.NoError(t, err)
	})

	return listChatsRepositoryFixture{
		repository:          repository,
		requesterID:         requesterID,
		peerID:              peerID,
		direct:              direct,
		group:               group,
		message:             message,
		peerDeletedAt:       peerDeletedAt,
		deletedPeerUsername: deletedPeerUsername,
	}
}

func insertListChatsUser(t *testing.T, pool *pgxpool.Pool, userID uuid.UUID) {
	t.Helper()
	_, err := pool.Exec(t.Context(), `
		INSERT INTO users (id, username, first_name, created_at, password_hash)
		VALUES ($1, $2, 'Chat List User', $3, 'password_hash')
	`, userID, listChatsUsername(userID), createDirectTestTime())
	require.NoError(t, err)
}

func insertListChatsGroup(
	t *testing.T,
	pool *pgxpool.Pool,
	group domain.GroupChat,
	participantIDs ...uuid.UUID,
) {
	t.Helper()
	tx, err := pool.Begin(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })

	_, err = tx.Exec(t.Context(), `
		INSERT INTO chats (id, type, last_message_id, last_activity_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, group.Chat.ID, group.Chat.Type, group.Chat.LastMessageID,
		group.Chat.LastActivityAt, group.Chat.CreatedAt)
	require.NoError(t, err)
	_, err = tx.Exec(t.Context(), `INSERT INTO groups (chat_id, title) VALUES ($1, $2)`,
		group.Chat.ID, group.Title)
	require.NoError(t, err)
	for index, participantID := range participantIDs {
		_, err = tx.Exec(t.Context(), `
			INSERT INTO chat_participants (chat_id, user_id, joined_at)
			VALUES ($1, $2, $3)
		`, group.Chat.ID, participantID, group.Chat.CreatedAt)
		require.NoError(t, err)
		role := "member"
		if index == 0 {
			role = "owner"
		}
		_, err = tx.Exec(t.Context(), `
			INSERT INTO group_participants (chat_id, user_id, role)
			VALUES ($1, $2, $3)
		`, group.Chat.ID, participantID, role)
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit(t.Context()))
}

func insertListChatsMessage(t *testing.T, pool *pgxpool.Pool, message domain.Message) {
	t.Helper()
	_, err := pool.Exec(t.Context(), `
		INSERT INTO messages (
			id, client_message_id, chat_id, sender_id, content, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, message.ID, message.ClientMessageID, message.ChatID, message.SenderID,
		message.Content, message.CreatedAt, message.UpdatedAt)
	require.NoError(t, err)
}

func listChatsUsername(userID uuid.UUID) string {
	return "list_" + strings.ReplaceAll(userID.String(), "-", "")[:16]
}
