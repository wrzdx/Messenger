package chats_service

import "errors"

var (
	ErrInvalidListChatsQuery = errors.New("invalid list chats query")
	ErrInvalidChatItem       = errors.New("invalid chat item")
)
