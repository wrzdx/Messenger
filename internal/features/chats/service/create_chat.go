package chats_service

import (
	"bytes"
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *ChatsService) CreateChat(
	ctx context.Context,
	userID uuid.UUID,
	chatType string,
	name *string,
	ParticipantIDs []uuid.UUID,
) (domain.Chat, error) {
	switch domain.ChatType(chatType) {
	case domain.ChatTypeDirect:
		return s.createDirect(ctx, userID, ParticipantIDs)

	case domain.ChatTypeGroup:
		return s.createGroup(ctx, userID, name, ParticipantIDs)

	default:
		return domain.Chat{}, domain.ErrInvalidChatType
	}
}

func (s *ChatsService) createGroup(
	ctx context.Context,
	userID uuid.UUID,
	name *string,
	ParticipantIDs []uuid.UUID,
) (domain.Chat, error) {

	// for _, id := range cmd.ParticipantIDs {
	// 	_, err := s.userRepository.GetByID(ctx, id)
	// 	if err != nil {
	// 		return domain.Chat{}, err
	// 	}
	// }

	// chat := domain.NewGroupChat(
	// 	uuid.New(),
	// 	*cmd.Title,
	// 	cmd.CreatorID,
	// )

	// if err := chat.Validate(); err != nil {
	// 	return domain.Chat{}, err
	// }

	// err := s.chatRepository.Create(ctx, chat, cmd.ParticipantIDs)
	// if err != nil {
	// 	return domain.Chat{}, err
	// }

	return domain.Chat{}, nil
}

func (s *ChatsService) createDirect(
	ctx context.Context,
	userID uuid.UUID,
	participantIDs []uuid.UUID,
) (domain.Chat, error) {
	if len(participantIDs) != 1 {
		return domain.Chat{}, domain.ErrInvalidParticipants
	}
	now := time.Now()
	chat := domain.NewChat(
		uuid.New(),
		domain.ChatTypeDirect,
		nil,
		nil,
		now,
		now,
	)
	other := participantIDs[0]
	user1, user2 := normalizeUsers(userID, other)
	chat, err := s.chatsRepo.CreateDirect(ctx, chat, user1, user2)
	if err != nil {
		return domain.Chat{}, fmt.Errorf("create chat: %w", err)
	}

	return chat, nil
}

func normalizeUsers(user1 uuid.UUID, user2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if bytes.Compare(user1[:], user2[:]) > 0 {
		user1, user2 = user2, user1
	}
	return user1, user2
}
