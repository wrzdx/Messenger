package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/go-chi/chi/v5"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
		credentials domain.UserCredentials,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	// GetUser(
	// 	ctx context.Context,
	// 	id int,
	// ) (domain.User, error)

	// DeleteUser(
	// 	ctx context.Context,
	// 	id int,
	// ) error

	// PatchUser(
	// 	ctx context.Context,
	// 	id int,
	// 	patch domain.UserPatch,
	// ) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/", h.CreateUser)
	router.Get("/", h.GetUsers)
	return router
}
