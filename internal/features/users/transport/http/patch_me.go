package users_transport_http

import (
	core_context "messenger/internal/core/context"
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
	claims := core_context.ClaimsRequired(ctx)
	sender := http_response.NewHTTPSender(log, w)

	var request PatchUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		sender.Error(err)
		return
	}

	userPatch := UserPatchFromRequest(request)

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

func (p *PatchUserRequest) Validate() map[string]string {
	fields := make(map[string]string)
	if p.Username.Set {
		if p.Username.Value == nil {
			fields["username"] = domain.ErrNullUsername.Error()
		} else if err := domain.ValidateUsername(*p.Username.Value); err != nil {
			fields["username"] = err.Error()
		}
	}

	if p.FirstName.Set {
		if p.FirstName.Value == nil {
			fields["first name"] = domain.ErrNullFirstname.Error()
		}
		if err := domain.ValidateFirstName(*p.FirstName.Value); err != nil {
			fields["first name"] = err.Error()
		}
	}
	if p.LastName.Set && p.LastName.Value != nil {
		if err := domain.ValidateLastName(*p.LastName.Value); err != nil {
			fields["last name"] = err.Error()
		}
	}

	if p.Bio.Set && p.Bio.Value != nil {
		if err := domain.ValidateBio(*p.Bio.Value); err != nil {
			fields["first name"] = err.Error()
		}
	}

	return fields
}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Username.ToDomain(),
		request.FirstName.ToDomain(),
		request.LastName.ToDomain(),
		request.Bio.ToDomain(),
	)
}
