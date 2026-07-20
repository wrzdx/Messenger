package chats_transport_http

import (
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type ChatResponse struct {
	ID             uuid.UUID  `json:"id"`
	Type           string     `json:"type"`
	LastMessageID  *uuid.UUID `json:"last_message_id"`
	LastActivityAt time.Time  `json:"last_activity_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func chatResponseFromDomain(chat domain.Chat) ChatResponse {
	return ChatResponse{
		ID:             chat.ID,
		Type:           string(chat.Type),
		LastMessageID:  chat.LastMessageID,
		LastActivityAt: chat.LastActivityAt,
		CreatedAt:      chat.CreatedAt,
	}
}
