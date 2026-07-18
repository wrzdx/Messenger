package chats_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type ChatsRepository interface {
	CreateDirect(
		ctx context.Context,
		direct domain.DirectChat,
		participant1, participant2 domain.ChatParticipant,
	) error

	GetDirectByUsers(
		ctx context.Context,
		user1ID, user2ID uuid.UUID,
	) (domain.DirectChat, error)
}

type TXManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type UsersRepository interface {
	GetUserForUpdate(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)
}
