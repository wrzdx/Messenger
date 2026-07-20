package chats_transport_http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	http_request "messenger/internal/core/transport/http/request"
	chats_service "messenger/internal/features/chats/service"
	"time"

	"github.com/google/uuid"
)

func encodeChatCursor(cursor *chats_service.ChatCursor) (*string, error) {
	if cursor == nil {
		return nil, nil
	}
	payload := chatCursorPayload{
		ChatID:         cursor.ChatID,
		LastActivityAt: cursor.LastActivityAt,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cursor: %w", err)
	}

	encoded := base64.URLEncoding.EncodeToString(jsonData)

	return &encoded, nil
}

func decodeChatCursor(cursorStr string) (*chats_service.ChatCursor, error) {
	var payload chatCursorPayload
	if cursorStr == "" {
		return nil, nil
	}

	jsonData, err := base64.URLEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode base64 cursor: %w",
			http_request.ErrInvalidRequest,
		)
	}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal cursor json: %w",
			http_request.ErrInvalidRequest,
		)
	}

	if payload.ChatID == uuid.Nil {
		return nil, fmt.Errorf(
			"cursor chat id is nil: %w",
			http_request.ErrInvalidRequest,
		)
	}

	if payload.LastActivityAt.IsZero() {
		return nil, fmt.Errorf(
			"cursor chat last_activity_at is zero value: %w",
			http_request.ErrInvalidRequest,
		)
	}

	return &chats_service.ChatCursor{
		ChatID:         payload.ChatID,
		LastActivityAt: payload.LastActivityAt,
	}, nil
}

type chatCursorPayload struct {
	ChatID         uuid.UUID `json:"chat_id"`
	LastActivityAt time.Time `json:"last_activity_at"`
}
