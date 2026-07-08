package http_request

import (
	"fmt"
	core_errors "messenger/internal/core/errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func GetQueryParam[T any](r *http.Request, key string) (*T, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	var target T

	switch any(&target).(type) {
	case *int:
		val, err := strconv.Atoi(param)
		if err != nil {
			return nil, fmt.Errorf("param='%s' by key='%s' not a valid integer: %w", param, key, core_errors.ErrValidation)
		}
		res := any(val).(T)
		return &res, nil

	case *uuid.UUID:
		id, err := uuid.Parse(param)
		if err != nil {
			return nil, fmt.Errorf("param='%s' by key='%s' not a valid uuid: %w", param, key, core_errors.ErrValidation)
		}
		res := any(id).(T)
		return &res, nil

	default:
		return nil, fmt.Errorf("unsupported type %T for query param conversion", target)
	}
}
