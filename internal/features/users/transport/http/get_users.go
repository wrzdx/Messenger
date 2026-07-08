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

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		err = fmt.Errorf(
			"%v: %w",
			err,
			http_response.ErrInvalidArgument,
		)
		responseHandler.ErrorResponse(
			http_response.Error{
				Error:   err,
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}
	pagination := domain.NewPagination(limit, offset)
	userDomains, err := h.usersService.GetUsers(ctx, pagination)
	if err != nil {
		responseHandler.ErrorResponse(http_response.MapError(err))
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
	limit, err := http_request.GetQueryParam[int](r, limitQueryParamKey)
	if err != nil {
		return nil, nil, err
	}

	offset, err := http_request.GetQueryParam[int](r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, err
	}

	return limit, offset, err
}
