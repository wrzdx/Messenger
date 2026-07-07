package users_transport_http

import (
	"fmt"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	core_http_types "messenger/internal/core/transport/http/types"
	"net/http"
)

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	claims, ok := core_auth.ClaimsFromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
	if !ok {
		responseHandler.ErrorResponse(
			core_http_response.MapError(core_http_response.ErrMissingClaims),
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			core_http_response.MapError(
				fmt.Errorf(
					"%v: %w",
					err,
					core_http_response.ErrInvalidArgument,
				),
			),
		)
		return
	}

	userPatch := UserPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, claims.UserID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

type PatchUserRequest struct {
	Username  core_http_types.Nullable[string] `json:"username" swaggertype:"string" example:"ivanov"`
	FirstName core_http_types.Nullable[string] `json:"first_name" swaggertype:"string" example:"Sidor"`
	LastName  core_http_types.Nullable[string] `json:"last_name" swaggertype:"string" example:"Ivanov"`
	Bio       core_http_types.Nullable[string] `json:"bio" swaggertype:"string" example:"I'like pizza!"`
}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Username.ToDomain(),
		request.FirstName.ToDomain(),
		request.LastName.ToDomain(),
		request.Bio.ToDomain(),
	)
}
