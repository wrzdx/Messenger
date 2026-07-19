package messages_transport_http

import (
	http_middleware "messenger/internal/core/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

type MessagesHandler struct {
	messagesService MessagesService
}

func NewMessagesHandler(messagesService MessagesService) *MessagesHandler {
	return &MessagesHandler{
		messagesService: messagesService,
	}
}

func (h *MessagesHandler) Router(authMW http_middleware.Middleware) chi.Router {
	router := chi.NewRouter()
	router.Use(authMW)
	router.Get("/", h.GetMessages)
	router.Post("/", h.SendMessage)
	return router
}
