package users_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	http_types "messenger/internal/core/transport/http/types"
	users_service "messenger/internal/features/users/service"
	"net/http"
)

func (h *UsersHandler) PatchMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	claims := core_context.ClaimsRequired(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)

	var request PatchProfileRequest
	if err := http_request.DecodeAndValidateRequestBody(r, &request); err != nil {
		sender.Error(err)
		return
	}
	user, err := h.usersService.UpdateProfile(ctx, claims.UserID, request.ToCommand())
	if err != nil {
		sender.Error(err)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	sender.OK(http.StatusOK, response)
}

type PatchProfileRequest struct {
	Username  *string                     `json:"username"`
	FirstName *string                     `json:"first_name"`
	LastName  http_types.Nullable[string] `json:"last_name"`
	Bio       http_types.Nullable[string] `json:"bio"`
}

func (r PatchProfileRequest) ToCommand() users_service.UpdateProfileCommand {
	return users_service.UpdateProfileCommand{
		Username:  r.Username,
		FirstName: r.FirstName,
		LastName:  r.LastName.ToCore(),
		Bio:       r.Bio.ToCore(),
	}
}
