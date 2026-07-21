package chats_transport_http

import (
	"fmt"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	chats_service "messenger/internal/features/chats/service"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (h *ChatsHandler) ListChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)
	claims := core_context.ClaimsRequired(ctx)

	queryParams := r.URL.Query()
	cursor, err := decodeChatCursor(queryParams.Get("cursor"))
	if err != nil {
		sender.Error(err)
		return
	}

	var limit int
	if limitString := queryParams.Get("limit"); limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			sender.Error(fmt.Errorf("invalid limit query param: %w", http_request.ErrInvalidRequest))
			return
		}
	}

	page, err := h.chatsService.ListChats(
		ctx,
		claims.UserID,
		chats_service.ListChatsQuery{Before: cursor, Limit: limit},
	)
	if err != nil {
		sender.Error(err)
		return
	}

	nextCursor, err := encodeChatCursor(page.NextCursor)
	if err != nil {
		sender.Error(err)
		return
	}

	response := ListChatsResponse{
		Chats:      make([]ChatItemResponse, len(page.Chats)),
		NextCursor: nextCursor,
	}
	for index, item := range page.Chats {
		response.Chats[index] = chatItemResponseFromService(item)
	}

	sender.OK(http.StatusOK, response)
}

type ListChatsResponse struct {
	Chats      []ChatItemResponse `json:"chats"`
	NextCursor *string            `json:"next_cursor"`
}

type ChatItemResponse struct {
	Chat        ChatResponse         `json:"chat"`
	Direct      *DirectChatResponse  `json:"direct,omitempty"`
	Group       *GroupChatResponse   `json:"group,omitempty"`
	LastMessage *LastMessageResponse `json:"last_message"`
}

type DirectChatResponse struct {
	Peer PeerResponse `json:"peer"`
}

type PeerResponse struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  *string    `json:"last_name"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type GroupChatResponse struct {
	Title string `json:"title"`
}

type LastMessageResponse struct {
	ID              uuid.UUID  `json:"id"`
	SenderID        uuid.UUID  `json:"sender_id"`
	SenderFirstName string     `json:"sender_first_name"`
	Content         string     `json:"content"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func chatItemResponseFromService(item chats_service.ChatItem) ChatItemResponse {
	response := ChatItemResponse{
		Chat: chatResponseFromDomain(item.Chat),
	}

	if item.DirectPeer != nil {
		response.Direct = &DirectChatResponse{
			Peer: PeerResponse{
				ID:        item.DirectPeer.ID,
				Username:  item.DirectPeer.Username,
				FirstName: item.DirectPeer.FirstName,
				LastName:  item.DirectPeer.LastName,
				DeletedAt: item.DirectPeer.DeletedAt,
			},
		}
	} else if item.GroupInfo != nil {
		response.Group = &GroupChatResponse{Title: item.GroupInfo.Title}
	}

	if item.LastMessage != nil {
		message := item.LastMessage.Message
		response.LastMessage = &LastMessageResponse{
			ID:              message.ID,
			SenderID:        message.SenderID,
			SenderFirstName: item.LastMessage.SenderFirstName,
			Content:         message.Content,
			CreatedAt:       message.CreatedAt,
			UpdatedAt:       message.UpdatedAt,
		}
	}
	return response
}
