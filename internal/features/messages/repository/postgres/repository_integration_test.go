//go:build integration

package messages_postgres_repository

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestGetMessageByClientID(t *testing.T) {
	config := postgres.NewConfigMust()
	pool, err := postgres.NewPool(t.Context(), config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	t.Run("distinguishes the same client id used by different senders", func(t *testing.T) {
		repository, chat, sender1ID := newMessageRepositoryFixture(t, pool, config.Timeout)
		sender2ID := uuid.New()
		insertMessageRepositoryUser(t, pool, sender2ID)
		t.Cleanup(func() {
			_, err := pool.Exec(context.Background(), `DELETE FROM messages WHERE sender_id = $1`, sender2ID)
			require.NoError(t, err)
			deleteMessageRepositoryUser(t, pool, config.Timeout, sender2ID)
		})
		clientMessageID := uuid.New()
		message1 := newRepositoryTestMessage(t, chat.ID, sender1ID, clientMessageID, "first")
		message2 := newRepositoryTestMessage(t, chat.ID, sender2ID, clientMessageID, "second")
		insertMessageRepositoryMessage(t, pool, message1)
		insertMessageRepositoryMessage(t, pool, message2)

		actual1, err := repository.GetMessageByClientID(t.Context(), sender1ID, clientMessageID)
		require.NoError(t, err)
		requireRepositoryMessageEqual(t, message1, actual1)

		actual2, err := repository.GetMessageByClientID(t.Context(), sender2ID, clientMessageID)
		require.NoError(t, err)
		requireRepositoryMessageEqual(t, message2, actual2)
	})

	t.Run("returns not found for unknown key", func(t *testing.T) {
		repository := NewRepository(pool, config.Timeout)

		message, err := repository.GetMessageByClientID(t.Context(), uuid.New(), uuid.New())

		require.ErrorIs(t, err, domain.ErrNotFound)
		require.Zero(t, message)
	})

	t.Run("rejects invalid state restored from database", func(t *testing.T) {
		repository, chat, senderID := newMessageRepositoryFixture(t, pool, config.Timeout)
		clientMessageID := uuid.New()
		_, err := pool.Exec(t.Context(), `
			INSERT INTO messages (
				id, client_message_id, chat_id, sender_id, content, created_at
			)
			VALUES ($1, $2, $3, $4, '', $5)
		`, uuid.New(), clientMessageID, chat.ID, senderID, repositoryTestTime())
		require.NoError(t, err)

		message, err := repository.GetMessageByClientID(t.Context(), senderID, clientMessageID)

		require.ErrorIs(t, err, domain.ErrInvalidMessage)
		require.Zero(t, message)
	})
}

func TestGetChatForUpdate(t *testing.T) {
	config := postgres.NewConfigMust()
	pool, err := postgres.NewPool(t.Context(), config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	t.Run("restores and validates chat", func(t *testing.T) {
		repository, expected, _ := newMessageRepositoryFixture(t, pool, config.Timeout)

		actual, err := repository.GetChatForUpdate(t.Context(), expected.ID)

		require.NoError(t, err)
		requireRepositoryChatEqual(t, expected, actual)
	})

	t.Run("returns not found for unknown chat", func(t *testing.T) {
		repository := NewRepository(pool, config.Timeout)

		chat, err := repository.GetChatForUpdate(t.Context(), uuid.New())

		require.ErrorIs(t, err, domain.ErrNotFound)
		require.Zero(t, chat)
	})

	t.Run("rejects invalid state restored from database", func(t *testing.T) {
		repository, chat, _ := newMessageRepositoryFixture(t, pool, config.Timeout)
		_, err := pool.Exec(t.Context(), `
			UPDATE chats
			SET last_activity_at = created_at - interval '1 second'
			WHERE id = $1
		`, chat.ID)
		require.NoError(t, err)

		actual, err := repository.GetChatForUpdate(t.Context(), chat.ID)

		require.ErrorIs(t, err, domain.ErrInvalidChat)
		require.Zero(t, actual)
	})

	t.Run("holds row lock until transaction ends", func(t *testing.T) {
		_, chat, _ := newMessageRepositoryFixture(t, pool, config.Timeout)
		tx, err := pool.Begin(t.Context())
		require.NoError(t, err)
		t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
		repository := NewRepository(tx, config.Timeout)

		_, err = repository.GetChatForUpdate(t.Context(), chat.ID)
		require.NoError(t, err)

		blockedCtx, cancel := context.WithTimeout(t.Context(), 150*time.Millisecond)
		defer cancel()
		_, err = pool.Exec(blockedCtx, `
			UPDATE chats SET last_activity_at = last_activity_at WHERE id = $1
		`, chat.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)

		require.NoError(t, tx.Rollback(t.Context()))
		_, err = pool.Exec(t.Context(), `
			UPDATE chats SET last_activity_at = last_activity_at WHERE id = $1
		`, chat.ID)
		require.NoError(t, err)
	})
}

func TestAppendMessage(t *testing.T) {
	config := postgres.NewConfigMust()
	pool, err := postgres.NewPool(t.Context(), config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	t.Run("persists message and updates chat state", func(t *testing.T) {
		repository, chat, senderID := newMessageRepositoryFixture(t, pool, config.Timeout)
		message := newRepositoryTestMessage(t, chat.ID, senderID, uuid.New(), "Hello")

		err := repository.AppendMessage(t.Context(), message)

		require.NoError(t, err)
		actual, err := repository.GetMessageByClientID(
			t.Context(),
			message.SenderID,
			message.ClientMessageID,
		)
		require.NoError(t, err)
		requireRepositoryMessageEqual(t, message, actual)
		requireMessageRepositoryChatState(t, pool, chat.ID, message.ID, message.CreatedAt)
	})

	t.Run("maps duplicate sender and client id", func(t *testing.T) {
		repository, chat, senderID := newMessageRepositoryFixture(t, pool, config.Timeout)
		clientMessageID := uuid.New()
		existing := newRepositoryTestMessage(t, chat.ID, senderID, clientMessageID, "first")
		require.NoError(t, repository.AppendMessage(t.Context(), existing))
		duplicate := newRepositoryTestMessage(t, chat.ID, senderID, clientMessageID, "second")

		err := repository.AppendMessage(t.Context(), duplicate)

		require.ErrorIs(t, err, domain.ErrAlreadyExists)
		actual, err := repository.GetMessageByClientID(t.Context(), senderID, clientMessageID)
		require.NoError(t, err)
		requireRepositoryMessageEqual(t, existing, actual)
		requireMessageRepositoryChatState(t, pool, chat.ID, existing.ID, existing.CreatedAt)
	})

	t.Run("rejects invalid message before persistence", func(t *testing.T) {
		repository := NewRepository(pool, config.Timeout)

		err := repository.AppendMessage(t.Context(), domain.Message{})

		require.ErrorIs(t, err, domain.ErrInvalidMessage)
	})

	t.Run("uses transaction from context and rolls back both changes", func(t *testing.T) {
		repository, chat, senderID := newMessageRepositoryFixture(t, pool, config.Timeout)
		message := newRepositoryTestMessage(t, chat.ID, senderID, uuid.New(), "rollback")
		manager := postgres.NewTransactionManager(pool)
		rollbackErr := errors.New("rollback test transaction")

		err := manager.WithinTransaction(t.Context(), func(ctx context.Context) error {
			if err := repository.AppendMessage(ctx, message); err != nil {
				return err
			}
			return rollbackErr
		})

		require.ErrorIs(t, err, rollbackErr)
		actual, err := repository.GetMessageByClientID(
			t.Context(),
			message.SenderID,
			message.ClientMessageID,
		)
		require.ErrorIs(t, err, domain.ErrNotFound)
		require.Zero(t, actual)
		requireMessageRepositoryInitialChatState(t, pool, chat)
	})
}

func newMessageRepositoryFixture(
	t *testing.T,
	pool *pgxpool.Pool,
	timeout time.Duration,
) (*Repository, domain.Chat, uuid.UUID) {
	t.Helper()
	userID := uuid.New()
	chat := domain.Chat{
		ID:             uuid.New(),
		Type:           domain.ChatTypeDirect,
		LastActivityAt: repositoryTestTime(),
		CreatedAt:      repositoryTestTime(),
	}
	chat.LastActivityAt = chat.CreatedAt
	require.NoError(t, chat.Validate())
	insertMessageRepositoryUser(t, pool, userID)
	_, err := pool.Exec(t.Context(), `
		INSERT INTO chats (id, type, last_message_id, last_activity_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, chat.ID, chat.Type, chat.LastMessageID, chat.LastActivityAt, chat.CreatedAt)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := pool.Exec(ctx, `DELETE FROM chats WHERE id = $1`, chat.ID)
		require.NoError(t, err)
		deleteMessageRepositoryUser(t, pool, timeout, userID)
	})
	return NewRepository(pool, timeout), chat, userID
}

func insertMessageRepositoryUser(t *testing.T, pool *pgxpool.Pool, userID uuid.UUID) {
	t.Helper()
	username := "msg_" + strings.ReplaceAll(userID.String(), "-", "")[:16]
	_, err := pool.Exec(t.Context(), `
		INSERT INTO users (id, username, first_name, created_at, password_hash)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, username, "Message Test", repositoryTestTime(), "password_hash")
	require.NoError(t, err)
}

func deleteMessageRepositoryUser(
	t *testing.T,
	pool *pgxpool.Pool,
	timeout time.Duration,
	userID uuid.UUID,
) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	require.NoError(t, err)
}

func insertMessageRepositoryMessage(
	t *testing.T,
	pool *pgxpool.Pool,
	message domain.Message,
) {
	t.Helper()
	_, err := pool.Exec(t.Context(), `
		INSERT INTO messages (
			id, client_message_id, chat_id, sender_id, content, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		message.ID,
		message.ClientMessageID,
		message.ChatID,
		message.SenderID,
		message.Content,
		message.CreatedAt,
		message.UpdatedAt,
	)
	require.NoError(t, err)
}

func newRepositoryTestMessage(
	t *testing.T,
	chatID, senderID, clientMessageID uuid.UUID,
	content string,
) domain.Message {
	t.Helper()
	message, err := domain.NewMessage(
		uuid.New(),
		clientMessageID,
		chatID,
		senderID,
		content,
		repositoryTestTime().Add(time.Minute),
	)
	require.NoError(t, err)
	return message
}

func requireRepositoryMessageEqual(t *testing.T, expected, actual domain.Message) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.ClientMessageID, actual.ClientMessageID)
	require.Equal(t, expected.ChatID, actual.ChatID)
	require.Equal(t, expected.SenderID, actual.SenderID)
	require.Equal(t, expected.Content, actual.Content)
	require.True(t, expected.CreatedAt.Equal(actual.CreatedAt))
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireRepositoryChatEqual(t *testing.T, expected, actual domain.Chat) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Type, actual.Type)
	require.Equal(t, expected.LastMessageID, actual.LastMessageID)
	require.True(t, expected.LastActivityAt.Equal(actual.LastActivityAt))
	require.True(t, expected.CreatedAt.Equal(actual.CreatedAt))
}

func requireMessageRepositoryChatState(
	t *testing.T,
	pool *pgxpool.Pool,
	chatID, messageID uuid.UUID,
	lastActivityAt time.Time,
) {
	t.Helper()
	var (
		actualMessageID      uuid.UUID
		actualLastActivityAt time.Time
	)
	err := pool.QueryRow(t.Context(), `
		SELECT last_message_id, last_activity_at FROM chats WHERE id = $1
	`, chatID).Scan(&actualMessageID, &actualLastActivityAt)
	require.NoError(t, err)
	require.Equal(t, messageID, actualMessageID)
	require.True(t, lastActivityAt.Equal(actualLastActivityAt))
}

func requireMessageRepositoryInitialChatState(
	t *testing.T,
	pool *pgxpool.Pool,
	chat domain.Chat,
) {
	t.Helper()
	var (
		lastMessageID *uuid.UUID
		lastActivity  time.Time
	)
	err := pool.QueryRow(t.Context(), `
		SELECT last_message_id, last_activity_at FROM chats WHERE id = $1
	`, chat.ID).Scan(&lastMessageID, &lastActivity)
	require.NoError(t, err)
	require.Nil(t, lastMessageID)
	require.True(t, chat.LastActivityAt.Equal(lastActivity))
}

func repositoryTestTime() time.Time {
	return time.Now().UTC().Truncate(time.Microsecond)
}
