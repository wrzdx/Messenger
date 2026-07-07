package users_transport_http

import (
	"fmt"
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
			http.StatusBadRequest,
			fmt.Errorf(
				"%w: %v",
				core_http_response.ErrInvalidArgument,
				err,
			),
		)
		return
	}

	if limit != nil && *limit < 0 {
		responseHandler.ErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf(
				"%w: limit must be non-negative",
				core_http_response.ErrInvalidArgument,
			),
		)
		return
	}

	if offset != nil && *offset < 0 {
		responseHandler.ErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf(
				"%w: offset must be non-negative",
				core_http_response.ErrInvalidArgument,
			),
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		statusCode := core_http_response.MapDomainErrorToStatusCode(err)
		responseHandler.ErrorResponse(statusCode, err)
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
