package users_transport_http

import (
	"fmt"
	auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	http_types "messenger/internal/core/transport/http/types"
	"net/http"
)

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	userID := auth.MustUserIDFromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request PatchUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			http_response.Error{
				Error: fmt.Errorf(
					"%v: %w",
					err,
					http_response.ErrInvalidArgument,
				),
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	userPatch := UserPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(http_response.MapError(err))
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

type PatchUserRequest struct {
	Username  http_types.Nullable[string] `json:"username" swaggertype:"string" example:"ivanov"`
	FirstName http_types.Nullable[string] `json:"first_name" swaggertype:"string" example:"Sidor"`
	LastName  http_types.Nullable[string] `json:"last_name" swaggertype:"string" example:"Ivanov"`
	Bio       http_types.Nullable[string] `json:"bio" swaggertype:"string" example:"I'like pizza!"`
}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Username.ToDomain(),
		request.FirstName.ToDomain(),
		request.LastName.ToDomain(),
		request.Bio.ToDomain(),
	)
}
