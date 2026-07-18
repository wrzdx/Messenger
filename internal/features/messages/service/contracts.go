package messages_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type MessagesRepository interface {
	GetMessageByClientID(
		ctx context.Context,
		senderID, clientMessageID uuid.UUID,
	) (domain.Message, error)

	AppendMessage(
		ctx context.Context,
		message domain.Message,
	) error
}

type ChatsRepository interface {
	GetChatForUpdate(
		ctx context.Context,
		chatID uuid.UUID,
	) (domain.Chat, error)

	GetDirectMessageState(
		ctx context.Context,
		chatID uuid.UUID,
	) (DirectMessageState, error)

	GetGroupSenderState(
		ctx context.Context,
		chatID, senderID uuid.UUID,
	) (GroupSenderState, error)
}

type TXManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
