package chats_transport_http

import (
	http_middleware "messenger/internal/core/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

type ChatsHandler struct {
	chatsService ChatsService
}

func NewChatsHandler(chatsService ChatsService) *ChatsHandler {
	return &ChatsHandler{
		chatsService: chatsService,
	}
}

func (h *ChatsHandler) Router(authMW http_middleware.Middleware) chi.Router {
	router := chi.NewRouter()
	router.Use(authMW)
	router.Get("/", h.ListChats)
	router.Post("/directs", h.CreateDirect)
	return router
}
