package users_transport_http

import (
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GetUserResponse UserDTOResponse

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)

	idStr := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		sender.Error(http_request.NewFieldError(map[string]string{
			"id": "invalid user uuid",
		}))
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		sender.Error(err)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	sender.OK(http.StatusOK, response)
}
