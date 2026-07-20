package chats_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s ChatsService) ListChats(
	ctx context.Context,
	requesterID uuid.UUID,
	query ListChatsQuery,
) (ChatPage, error) {
	if requesterID == uuid.Nil {
		return ChatPage{}, fmt.Errorf("nil requester id: %w", domain.ErrNotFound)
	}

	query = query.normalize()
	if err := query.validate(); err != nil {
		return ChatPage{}, fmt.Errorf("validate query: %w", err)
	}

	allChats, err := s.chatsRepo.ListChats(
		ctx,
		requesterID,
		query.Before,
		query.Limit+1,
	)

	if err != nil {
		return ChatPage{}, fmt.Errorf("list chats: %w", err)
	}
	if allChats == nil {
		allChats = make([]ChatItem, 0)
	}

	var page ChatPage
	hasMore := len(allChats) > query.Limit
	if hasMore {
		allChats = allChats[:query.Limit]
	}

	page.Chats = allChats

	if hasMore {
		last := page.Chats[len(page.Chats)-1]
		chat, err := last.baseChat()
		if err != nil {
			return ChatPage{}, fmt.Errorf("get last chat: %w", err)
		}
		page.NextCursor = &ChatCursor{
			ChatID:         chat.ID,
			LastActivityAt: chat.LastActivityAt,
		}
	}
	return page, nil
}

type ListChatsQuery struct {
	Before *ChatCursor
	Limit  int
}

func (q ListChatsQuery) normalize() ListChatsQuery {
	if q.Limit == 0 {
		q.Limit = 50
	}
	return q
}

func (q ListChatsQuery) validate() error {
	fields := make(map[string]string)
	if q.Limit < 0 || q.Limit > 100 {
		fields["limit"] = "limit must be between 1 and 100"
	}

	if q.Before != nil {
		if q.Before.LastActivityAt.IsZero() {
			fields["last_activity_at"] = "last_activity_at of chat cursor cannot be zero value"
		}
		if q.Before.ChatID == uuid.Nil {
			fields["chat_id"] = "chat id of chat cursor cannot be nil"
		}
	}
	if len(fields) > 0 {
		return domain.DetailedError{
			Err:     ErrInvalidListChatsQuery,
			Details: fields,
		}
	}

	return nil
}

type ChatItem struct {
	Direct      *DirectChatItem
	Group       *domain.GroupChat
	LastMessage *LastMessageItem
}

func (i ChatItem) Validate() error {
	chat, err := i.baseChat()
	if err != nil {
		return err
	}

	if i.Direct != nil {
		if err := i.Direct.Chat.Validate(); err != nil {
			return fmt.Errorf("validate direct chat: %v: %w", err, ErrInvalidChatItem)
		}
		if i.Direct.PeerID != i.Direct.Chat.User1ID &&
			i.Direct.PeerID != i.Direct.Chat.User2ID {
			return fmt.Errorf("peer does not belong to direct chat: %w", ErrInvalidChatItem)
		}
		if err := i.Direct.PeerProfile.Validate(); err != nil {
			return fmt.Errorf("validate peer profile: %v: %w", err, ErrInvalidChatItem)
		}
		if i.Direct.PeerDeletedAt != nil && i.Direct.PeerDeletedAt.IsZero() {
			return fmt.Errorf("peer deleted_at is zero: %w", ErrInvalidChatItem)
		}
	}

	if i.Group != nil {
		if err := i.Group.Validate(); err != nil {
			return fmt.Errorf("validate group chat: %v: %w", err, ErrInvalidChatItem)
		}
	}

	if chat.LastMessageID == nil {
		if i.LastMessage != nil {
			return fmt.Errorf("unexpected last message: %w", ErrInvalidChatItem)
		}
		return nil
	}

	if i.LastMessage == nil {
		return fmt.Errorf("last message is missing: %w", ErrInvalidChatItem)
	}
	if err := i.LastMessage.Message.Validate(); err != nil {
		return fmt.Errorf("validate last message: %v: %w", err, ErrInvalidChatItem)
	}
	if i.LastMessage.Message.ID != *chat.LastMessageID {
		return fmt.Errorf("last message id mismatch: %w", ErrInvalidChatItem)
	}
	if i.LastMessage.Message.ChatID != chat.ID {
		return fmt.Errorf("last message chat id mismatch: %w", ErrInvalidChatItem)
	}
	if err := i.LastMessage.SenderProfile.Validate(); err != nil {
		return fmt.Errorf("validate last message sender profile: %v: %w", err, ErrInvalidChatItem)
	}

	return nil
}

type DirectChatItem struct {
	Chat          domain.DirectChat
	PeerID        uuid.UUID
	PeerProfile   domain.UserProfile
	PeerDeletedAt *time.Time
}

type LastMessageItem struct {
	Message       domain.Message
	SenderProfile domain.UserProfile
}

func (i ChatItem) baseChat() (domain.Chat, error) {
	switch {
	case i.Direct != nil && i.Group == nil:
		return i.Direct.Chat.Chat, nil
	case i.Direct == nil && i.Group != nil:
		return i.Group.Chat, nil
	default:
		return domain.Chat{}, fmt.Errorf(
			"chat item must contain exactly one chat subtype: %w",
			ErrInvalidChatItem,
		)
	}
}

type ChatPage struct {
	Chats      []ChatItem
	NextCursor *ChatCursor
}

type ChatCursor struct {
	ChatID         uuid.UUID
	LastActivityAt time.Time
}
