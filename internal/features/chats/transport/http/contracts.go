package chats_transport_http

import (
	"context"
	"messenger/internal/core/domain"
	chats_service "messenger/internal/features/chats/service"

	"github.com/google/uuid"
)

type ChatsService interface {
	CreateDirect(
		ctx context.Context,
		currentUserID uuid.UUID,
		peerID uuid.UUID,
	) (domain.DirectChat, bool, error)

	ListChats(
		ctx context.Context,
		requesterID uuid.UUID,
		query chats_service.ListChatsQuery,
	) (chats_service.ChatPage, error)
}
