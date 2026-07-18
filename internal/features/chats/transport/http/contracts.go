package chats_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type ChatsService interface {
	CreateDirect(
		ctx context.Context,
		currentUserID uuid.UUID,
		peerID uuid.UUID,
	) (domain.DirectChat, bool, error)
}
