package users_transport_http

import (
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func ValidatePagination(limit, offset *int) map[string]string {
	fields := make(map[string]string)
	if limit != nil {
		if err := domain.ValidateLimit(*limit); err != nil {
			fields["limit"] = err.Error()
		}
	}

	if offset != nil {
		if err := domain.ValidateOffset(*offset); err != nil {
			fields["offset"] = err.Error()
		}
	}

	return fields
}

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		sender.Error(err)
		return
	}
	if fields := ValidatePagination(limit, offset); len(fields) > 0 {
		sender.Error(core_errors.ValidationError(fields))
		return
	}
	pagination := domain.NewPagination(limit, offset)
	userDomains, err := h.usersService.GetUsers(ctx, pagination)
	if err != nil {
		sender.Error(err)
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	sender.OK(http.StatusOK, response)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := http_request.GetQueryParam[int](r, limitQueryParamKey)
	if err != nil {
		return nil, nil, err
	}

	offset, err := http_request.GetQueryParam[int](r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, err
	}

	return limit, offset, nil
}
