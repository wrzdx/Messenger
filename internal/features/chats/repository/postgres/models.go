package chats_postgres_repository

import (
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type ChatModel struct {
	ID             uuid.UUID
	Type           domain.ChatType
	Title          *string
	LastMessageID  *uuid.UUID
	LastActivityAt time.Time
	CreatedAt      time.Time
}
