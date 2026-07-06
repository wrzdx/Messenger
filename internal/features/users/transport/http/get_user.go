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
		responseHandler.ErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf("%w: %v", ErrInvalidArgument, err),
		)
		return
	}
	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		statusCode := mapDomainErrorToStatusCode(err)
		responseHandler.ErrorResponse(statusCode, err)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
