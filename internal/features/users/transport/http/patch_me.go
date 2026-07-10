package users_transport_http

import (
	"fmt"
	core_context "messenger/internal/core/context"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	http_types "messenger/internal/core/transport/http/types"
	users_service "messenger/internal/features/users/service"
	"net/http"
)

type PatchUserResponse UserDTOResponse

func (h *UsersHandler) PatchMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	claims := core_context.ClaimsRequired(ctx)
	sender := http_response.NewHTTPSender(log, w)

	var request PatchUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		fmt.Printf("%+v\n", err)
		sender.Error(err)
		return
	}

	userPatch := users_service.UserPatch{
		Username:  request.Username,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Bio:       request.Bio,
	}

	userDomain, err := h.usersService.PatchUser(ctx, claims.UserID, userPatch)
	if err != nil {
		sender.Error(err)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	sender.OK(http.StatusOK, response)
}

type PatchUserRequest struct {
	Username  http_types.Nullable[string] `json:"username" swaggertype:"string" example:"ivanov"`
	FirstName http_types.Nullable[string] `json:"first_name" swaggertype:"string" example:"Sidor"`
	LastName  http_types.Nullable[string] `json:"last_name" swaggertype:"string" example:"Ivanov"`
	Bio       http_types.Nullable[string] `json:"bio" swaggertype:"string" example:"I'like pizza!"`
}
