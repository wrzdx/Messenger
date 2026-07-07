package users_transport_http

import (
	"fmt"
	"messenger/internal/core/domain"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
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
	pagination := domain.NewPagination(limit, offset)
	userDomains, err := h.usersService.GetUsers(ctx, pagination)
	if err != nil {
		responseHandler.ErrorResponse(core_http_response.MapError(err))
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_request.GetQueryParam[int](r, limitQueryParamKey)
	if err != nil {
		return nil, nil, err
	}

	offset, err := core_http_request.GetQueryParam[int](r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, err
	}

	return limit, offset, err
}
