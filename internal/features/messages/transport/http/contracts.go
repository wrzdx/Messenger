package messages_transport_http

import (
	"context"
	"messenger/internal/core/domain"
	messages_service "messenger/internal/features/messages/service"

	"github.com/google/uuid"
)

type MessagesService interface {
	GetMessages(
		ctx context.Context,
		requesterID uuid.UUID,
		query messages_service.GetMessagesQuery,
	) (messages_service.MessagePage, error)

	SendMessage(
		ctx context.Context,
		senderID uuid.UUID,
		command messages_service.SendMessageCommand,
	) (domain.Message, bool, error)
}
