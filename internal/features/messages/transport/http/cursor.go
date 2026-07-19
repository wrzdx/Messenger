package messages_transport_http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	http_request "messenger/internal/core/transport/http/request"
	messages_service "messenger/internal/features/messages/service"
	"time"

	"github.com/google/uuid"
)

func encodeMessageCursor(cursor *messages_service.MessageCursor) (*string, error) {
	if cursor == nil {
		return nil, nil
	}
	payload := messageCursorPayload{
		MessageID: cursor.MessageID,
		CreatedAt: cursor.CreatedAt,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cursor: %w", err)
	}

	encoded := base64.URLEncoding.EncodeToString(jsonData)

	return &encoded, nil
}

func decodeMessageCursor(cursorStr string) (*messages_service.MessageCursor, error) {
	var payload messageCursorPayload
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

	if payload.MessageID == uuid.Nil {
		return nil, fmt.Errorf(
			"cursor message id is nil: %w",
			http_request.ErrInvalidRequest,
		)
	}

	if payload.CreatedAt.IsZero() {
		return nil, fmt.Errorf(
			"cursor message created_at is zero value: %w",
			http_request.ErrInvalidRequest,
		)
	}

	return &messages_service.MessageCursor{
		MessageID: payload.MessageID,
		CreatedAt: payload.CreatedAt,
	}, nil
}

type messageCursorPayload struct {
	MessageID uuid.UUID `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}
