package chats_service

import "errors"

var (
	ErrInvalidListChatsQuery             = errors.New("invalid list chats query")
	ErrInvalidChatItem                   = errors.New("invalid chat item")
	ErrInvalidListGroupParticipantsQuery = errors.New("invalid group participant query")
	ErrInvalidGroupParticipantItem       = errors.New("invalid group participant item")
)
