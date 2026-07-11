package chats_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/repository/postgres"

	"github.com/google/uuid"
)

func (r *ChatsRepository) CreateDirect(
	ctx context.Context,
	chat domain.Chat,
	user1 uuid.UUID,
	user2 uuid.UUID,
) error {
	context, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()

	query := `
INSERT INTO chats (id, type, last_activity_at, created_at)
VALUES ($1, $2,$3,$4) 
	RETURNING 
		id,
		type,
		title,
		last_message_id,
		last_activity_at,
		created_at,
	`

	var model ChatModel
	err := r.db.QueryRow(ctx, query, chat.ID, chat.Type, chat.LastActivityAt, chat.CreatedAt).
		Scan(&model.ID, &model.Type, &model.Title, &model.LastMessageID, &model.LastActivityAt, &model.CreatedAt)

	if err != nil {
		return fmt.Errorf("create chat: %w", err)
	}

	query = `
	INSERT INTO directs (chat_id, user1_id, user_2_id)
	VALUES ($1, $2, $3);
	`

	_, err := r.db.Exec(ctx, query, chat.ID, user1, user2)
	if err != nil {
		if errors.Is(err, postgres.ErrViolatesUnique) {
			failedField, failedValue := getConstraintValues(user, err)
			return domain.User{}, domain.AlreadyExistsErr(
				domain.UserEntity,
				failedField,
				failedValue,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}
}
