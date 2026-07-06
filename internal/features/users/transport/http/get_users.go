package users_transport_http

import (
	"fmt"
	core_errors "messenger/internal/core/errors"
	core_logger "messenger/internal/core/logger"
	core_http_request "messenger/internal/core/transport/http/request"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

type GetUsersResponse []UserDTOResponse

// GetUsers      godoc
// @Summary     Список пользователей
// @Description Просмотр списка пользователей с опциональной пагинацией
// @Tags        users
// @Produce     json
// @Param       limit query int false "Размер страницы с пользователями"
// @Param       offset query int false "Смещение страницы с пользователями"
// @Success     200 {object} GetUsersResponse "Успешное получение списка пользователей"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router      /users [get]
func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit/offset' query param")
		return
	}

	if limit != nil && *limit < 0 {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "limit must be non-negative")
		return
	}

	if offset != nil && *offset < 0 {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "offset must be non-negative")
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
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
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, err
}
