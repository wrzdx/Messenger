package chats_transport_http

import "github.com/go-chi/chi/v5"

type ChatsHandler struct {
	chatsService ChatsService
}

type ChatsService interface {
}

func NewChatsHandler(chatsService ChatsService) ChatsHandler {
	return ChatsHandler{
		chatsService: chatsService,
	}
}

func (h *ChatsHandler) Router() chi.Router {
	router := chi.NewRouter()
	return router
}
