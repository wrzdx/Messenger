package users_transport_http

import (
	"fmt"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"

	"github.com/google/uuid"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_request.GetPathValue[uuid.UUID](r, "id")
	if err != nil {
		err = fmt.Errorf(
			"%v: %w",
			err,
			core_http_response.ErrInvalidArgument,
		)
		responseHandler.ErrorResponse(
			core_http_response.Error{
				Error:   err,
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}
	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
