package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"
	http_middleware "messenger/internal/core/transport/http/middleware"

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

func (h *UsersHTTPHandler) Router(authMW http_middleware.Middleware) chi.Router {
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
