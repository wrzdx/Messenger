package realtime_transport_ws

import (
	"time"

	"github.com/google/uuid"
)

const (
	authReqType = "authenticate"
	authResType = "authenticated"
	msgCreated = "message_created"
)

type AuthRequest struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	Type string `json:"type"`
}

// Event is the server-to-client envelope. Data holds a transport DTO and is
// serialized once by Publish before being shared across connection queues.
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

type Message struct {
	ID              uuid.UUID  `json:"id"`
	ChatID          uuid.UUID  `json:"chat_id"`
	SenderID        uuid.UUID  `json:"sender_id"`
	Content         string     `json:"content"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
