package users_transport_http

import (
	"fmt"
	"messenger/internal/core/domain"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type CreateUserRequest struct {
	Username  string  `json:"username" validate:"required,min=5,max=32" example:"qwerty"`
	FirstName string  `json:"first_name" validate:"required,min=1,max=64" example:"Ivan"`
	LastName  *string `json:"last_name" validate:"omitempty,max=64" example:"Ivanov"`
	Bio       *string `json:"bio" validate:"omitempty,max=70" example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
	Password  string  `json:"password" validate:"required,min=8,max=32" example:"password"`
}

type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf("%w: %v", ErrInvalidArgument, err),
		)
		return
	}
	payload := domain.NewRegisterUserPayload(
		request.Username,
		request.FirstName,
		request.LastName,
		request.Bio,
		request.Password,
	)
	userDomain, err := h.usersService.CreateUser(ctx, payload)
	if err != nil {
		statusCode := mapDomainErrorToStatusCode(err)
		responseHandler.ErrorResponse(statusCode, err)
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
