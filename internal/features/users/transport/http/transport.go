package users_transport_http

import (
	http_middleware "messenger/internal/core/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	usersService UsersService
}

func NewUsersHandler(usersService UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (h *UsersHandler) Router(authMW http_middleware.Middleware) chi.Router {
	router := chi.NewRouter()
	router.Get("/{id}", h.GetUser)
	router.Get("/me", h.GetMe)

	return router
}
