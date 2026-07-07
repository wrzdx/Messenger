package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	GetUsers(
		ctx context.Context,
		pagination domain.Pagination,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		patch domain.UserPatch,
	) (domain.User, error)

	ChangePassword(
		ctx context.Context,
		id uuid.UUID,
		old_password string,
		new_password string,
	) error
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/", h.GetUsers)
	router.Get("/me", h.GetMe)
	router.Get("/{id}", h.GetUser)
	router.Patch("/me", h.PatchMe)
	router.Delete("/me", h.DeleteMe)
	router.Put("/me/password", h.ChangePassword)
	return router
}
