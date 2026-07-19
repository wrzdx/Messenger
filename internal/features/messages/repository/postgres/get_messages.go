package messages_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"
	messages_service "messenger/internal/features/messages/service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetMessages(
	ctx context.Context,
	chatID uuid.UUID,
	before *messages_service.MessageCursor,
	limit int,
) ([]domain.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	db := postgres.GetExecutor(ctx, r.db)

	nonspecifiedQuery := `
	SELECT id, client_message_id, chat_id, sender_id, content, created_at, updated_at
	FROM messages
	WHERE chat_id = $1
	ORDER BY created_at DESC, id DESC
	LIMIT $2;
	`
	specifiedQuery := `
	SELECT id, client_message_id, chat_id, sender_id, content, created_at, updated_at
	FROM messages
	WHERE chat_id = $1
	AND (created_at, id) < ($2, $3)
	ORDER BY created_at DESC, id DESC
	LIMIT $4;
	`

	var (
		messages []domain.Message
		rows     pgx.Rows
		err      error
	)

	if before != nil {
		rows, err = db.Query(
			ctx,
			specifiedQuery,
			chatID,
			before.CreatedAt,
			before.MessageID,
			limit,
		)
	} else {
		rows, err = db.Query(
			ctx,
			nonspecifiedQuery,
			chatID,
			limit,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m domain.Message
		err := rows.Scan(
			&m.ID,
			&m.ClientMessageID,
			&m.ChatID,
			&m.SenderID,
			&m.Content,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		if err := m.Validate(); err != nil {
			return nil, fmt.Errorf("validate message from db: %w", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}
