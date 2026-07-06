package users_transport_http

import (
	"messenger/internal/core/domain"

	"github.com/go-chi/chi/v5"
)

type UsersHTTPHandler struct {
	usersService domain.UsersService
}

func NewUsersHTTPHandler(usersService domain.UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/", h.CreateUser)
	router.Get("/", h.GetUsers)
	router.Get("/me", h.GetMe)
	router.Delete("/me", h.DeleteMe)
	router.Patch("/me", h.PatchMe)
	router.Get("/{id}", h.GetUser)
	return router
}
