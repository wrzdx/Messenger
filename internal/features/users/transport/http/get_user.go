package users_transport_http

import (
	"fmt"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"

	"github.com/google/uuid"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userID, err := http_request.GetPathValue[uuid.UUID](r, "id")
	if err != nil {
		err = fmt.Errorf(
			"%v: %w",
			err,
			http_response.ErrInvalidArgument,
		)
		responseHandler.ErrorResponse(
			http_response.Error{
				Error:   err,
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}
	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(http_response.MapError(err))
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
