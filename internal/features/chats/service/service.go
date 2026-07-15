package chats_service

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type ChatsService struct {
	txmanager TXManager
	chatsRepo ChatsRepository
	usersRepo UsersRepostiory
}

type ChatsRepository interface {
	CreateDirect(
		ctx context.Context,
		chat domain.Chat,
		user1 uuid.UUID,
		user2 uuid.UUID,
	) (domain.Chat, error)
}

type TXManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type UsersRepostiory interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)
}

func NewChatsService(repo ChatsRepository, txmanager TXManager) *ChatsService {
	return &ChatsService{
		chatsRepo: repo,
	}
}
