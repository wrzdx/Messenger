package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"
	http_middleware "messenger/internal/core/transport/http/middleware"
	users_service "messenger/internal/features/users/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsersHandler struct {
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
		patch users_service.UserPatch,
	) (domain.User, error)
}

func NewUsersHandler(usersService UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (h *UsersHandler) Router(authMW http_middleware.Middleware) chi.Router {
	router := chi.NewRouter()
	router.Get("/", h.GetUsers)
	router.Get("/{id}", h.GetUser)
	router.Group(func(r chi.Router) {
		r.Use(authMW)
		r.Get("/me", h.GetMe)
		r.Patch("/me", h.PatchMe)
		r.Delete("/me", h.DeleteMe)
		r.Put("/me/password", h.ChangePassword)
	})

	return router
}
