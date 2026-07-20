package chats_postgres_repository

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/postgres"
	chats_service "messenger/internal/features/chats/service"
	"time"

	"github.com/google/uuid"
)

func (r *ChatsRepository) ListChats(
	ctx context.Context,
	userID uuid.UUID,
	before *chats_service.ChatCursor,
	limit int,
) ([]chats_service.ChatItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	db := postgres.GetExecutor(ctx, r.db)

	query := `
	WITH page AS (
		SELECT c.id, c.type, c.last_message_id, c.last_activity_at, c.created_at
		FROM chat_participants cp
		JOIN chats c ON c.id = cp.chat_id
		WHERE cp.user_id = $1
		  AND (
			$2::timestamptz IS NULL
			OR (c.last_activity_at, c.id) < ($2, $3)
		  )
		ORDER BY c.last_activity_at DESC, c.id DESC
		LIMIT $4
	)
	SELECT c.id, c.type, c.last_message_id, c.last_activity_at, c.created_at,
	       d.user1_id, d.user2_id,
	       peer.id, peer.username, peer.first_name, peer.last_name, peer.bio, peer.deleted_at,
	       g.title,
	       lm.id, lm.client_message_id, lm.chat_id, lm.sender_id,
	       lm.content, lm.created_at, lm.updated_at,
	       message_sender.username, message_sender.first_name,
	       message_sender.last_name, message_sender.bio
	FROM page c
	LEFT JOIN directs d ON d.chat_id = c.id
	LEFT JOIN groups g ON g.chat_id = c.id
	LEFT JOIN users peer ON peer.id = CASE
		WHEN d.user1_id = $1 THEN d.user2_id
		ELSE d.user1_id
	END
	LEFT JOIN messages lm ON lm.id = c.last_message_id
	LEFT JOIN users message_sender ON message_sender.id = lm.sender_id
	ORDER BY c.last_activity_at DESC, c.id DESC;
	`

	var beforeActivity *time.Time
	var beforeChatID *uuid.UUID
	if before != nil {
		activity := before.LastActivityAt
		chatID := before.ChatID
		beforeActivity = &activity
		beforeChatID = &chatID
	}

	rows, err := db.Query(ctx, query, userID, beforeActivity, beforeChatID, limit)
	if err != nil {
		return nil, fmt.Errorf("list chats: %w", err)
	}
	defer rows.Close()

	var chats []chats_service.ChatItem
	for rows.Next() {
		var row listChatRow
		if err := rows.Scan(
			&row.chat.ID,
			&row.chat.Type,
			&row.chat.LastMessageID,
			&row.chat.LastActivityAt,
			&row.chat.CreatedAt,
			&row.directUser1ID,
			&row.directUser2ID,
			&row.peer.id,
			&row.peer.username,
			&row.peer.firstName,
			&row.peer.lastName,
			&row.peer.bio,
			&row.peer.deletedAt,
			&row.groupTitle,
			&row.message.id,
			&row.message.clientMessageID,
			&row.message.chatID,
			&row.message.senderID,
			&row.message.content,
			&row.message.createdAt,
			&row.message.updatedAt,
			&row.sender.username,
			&row.sender.firstName,
			&row.sender.lastName,
			&row.sender.bio,
		); err != nil {
			return nil, fmt.Errorf("scan chat list item: %w", err)
		}

		item, err := row.toService(userID)
		if err != nil {
			return nil, fmt.Errorf("restore chat list item: %w", err)
		}
		chats = append(chats, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chat list: %w", err)
	}

	return chats, nil
}

type listChatRow struct {
	chat domain.Chat

	directUser1ID *uuid.UUID
	directUser2ID *uuid.UUID
	peer          nullableUserProfile
	groupTitle    *string

	message nullableMessage
	sender  nullableUserProfile
}

type nullableUserProfile struct {
	id        *uuid.UUID
	username  *string
	firstName *string
	lastName  *string
	bio       *string
	deletedAt *time.Time
}

type nullableMessage struct {
	id              *uuid.UUID
	clientMessageID *uuid.UUID
	chatID          *uuid.UUID
	senderID        *uuid.UUID
	content         *string
	createdAt       *time.Time
	updatedAt       *time.Time
}

func (r listChatRow) toService(requesterID uuid.UUID) (chats_service.ChatItem, error) {
	if err := r.chat.Validate(); err != nil {
		return chats_service.ChatItem{}, fmt.Errorf("validate chat: %w", err)
	}

	item := chats_service.ChatItem{}
	switch r.chat.Type {
	case domain.ChatTypeDirect:
		if r.directUser1ID == nil || r.directUser2ID == nil ||
			r.peer.id == nil || r.peer.username == nil || r.peer.firstName == nil {
			return chats_service.ChatItem{}, fmt.Errorf("incomplete direct chat row")
		}
		if r.groupTitle != nil {
			return chats_service.ChatItem{}, fmt.Errorf("direct chat has group data")
		}
		if requesterID != *r.directUser1ID && requesterID != *r.directUser2ID {
			return chats_service.ChatItem{}, fmt.Errorf("requester does not belong to direct chat")
		}
		if *r.peer.id == requesterID {
			return chats_service.ChatItem{}, fmt.Errorf("direct peer equals requester")
		}

		profile, err := domain.NewUserProfile(
			*r.peer.username,
			*r.peer.firstName,
			r.peer.lastName,
			r.peer.bio,
		)
		if err != nil {
			return chats_service.ChatItem{}, fmt.Errorf("restore direct peer profile: %w", err)
		}

		direct := domain.DirectChat{
			Chat:    r.chat,
			User1ID: *r.directUser1ID,
			User2ID: *r.directUser2ID,
		}
		if err := direct.Validate(); err != nil {
			return chats_service.ChatItem{}, fmt.Errorf("validate direct chat: %w", err)
		}
		item.Direct = &chats_service.DirectChatItem{
			Chat:          direct,
			PeerID:        *r.peer.id,
			PeerProfile:   profile,
			PeerDeletedAt: r.peer.deletedAt,
		}

	case domain.ChatTypeGroup:
		if r.groupTitle == nil {
			return chats_service.ChatItem{}, fmt.Errorf("group title is missing")
		}
		if r.directUser1ID != nil || r.directUser2ID != nil || r.peer.id != nil {
			return chats_service.ChatItem{}, fmt.Errorf("group chat has direct data")
		}
		group := domain.GroupChat{Chat: r.chat, Title: *r.groupTitle}
		if err := group.Validate(); err != nil {
			return chats_service.ChatItem{}, fmt.Errorf("validate group chat: %w", err)
		}
		item.Group = &group

	default:
		return chats_service.ChatItem{}, fmt.Errorf("unsupported chat type %q", r.chat.Type)
	}

	if r.chat.LastMessageID != nil {
		if r.message.id == nil || r.message.clientMessageID == nil ||
			r.message.chatID == nil || r.message.senderID == nil ||
			r.message.content == nil || r.message.createdAt == nil ||
			r.sender.username == nil || r.sender.firstName == nil {
			return chats_service.ChatItem{}, fmt.Errorf("incomplete last message row")
		}

		message := domain.Message{
			ID:              *r.message.id,
			ClientMessageID: *r.message.clientMessageID,
			ChatID:          *r.message.chatID,
			SenderID:        *r.message.senderID,
			Content:         *r.message.content,
			CreatedAt:       *r.message.createdAt,
			UpdatedAt:       r.message.updatedAt,
		}
		if err := message.Validate(); err != nil {
			return chats_service.ChatItem{}, fmt.Errorf("validate last message: %w", err)
		}

		senderProfile, err := domain.NewUserProfile(
			*r.sender.username,
			*r.sender.firstName,
			r.sender.lastName,
			r.sender.bio,
		)
		if err != nil {
			return chats_service.ChatItem{}, fmt.Errorf("restore message sender profile: %w", err)
		}
		item.LastMessage = &chats_service.LastMessageItem{
			Message:       message,
			SenderProfile: senderProfile,
		}
	} else if r.message.id != nil {
		return chats_service.ChatItem{}, fmt.Errorf("unexpected last message row")
	}

	if err := item.Validate(); err != nil {
		return chats_service.ChatItem{}, fmt.Errorf("validate chat item: %w", err)
	}
	return item, nil
}
