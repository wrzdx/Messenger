package chats_service

import "errors"

var (
	ErrInvalidInput                = errors.New("invalid input")
	ErrInvalidChatItem             = errors.New("invalid chat item")
	ErrInvalidGroupParticipantItem = errors.New("invalid group participant item")
	ErrNotEnoughRights             = errors.New("not enough rights")
	ErrOwnerCannotQuitGroup        = errors.New("owner cannot quit group")
)
