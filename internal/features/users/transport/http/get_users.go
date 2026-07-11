package users_transport_http

import (
	"fmt"
	"messenger/internal/core/domain"
	logger "messenger/internal/core/logger"
	http_request "messenger/internal/core/transport/http/request"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	sender := http_response.NewHTTPSender(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		sender.Error(fmt.Errorf("%w: %v", domain.ErrValidation, err))
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
