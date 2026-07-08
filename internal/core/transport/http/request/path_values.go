package http_request

import (
	"fmt"
	core_errors "messenger/internal/core/errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func GetPathValue[T any](r *http.Request, key string) (T, error) {
	var zero T
	pathValue := r.PathValue(key)

	if pathValue == "" {
		return zero, fmt.Errorf("no key='%s' in path values", key)
	}

	switch any(&zero).(type) {
	case *int:
		val, err := strconv.Atoi(pathValue)
		if err != nil {
			return zero, fmt.Errorf("path value='%s' by key='%s' not a valid integer: %w", pathValue, key, core_errors.ErrValidation)
		}
		return any(val).(T), nil

	case *uuid.UUID:
		val, err := uuid.Parse(pathValue)
		if err != nil {
			return zero, fmt.Errorf("path value='%s' by key='%s' not a valid uuid: %w", pathValue, key,  core_errors.ErrValidation)
		}
		return any(val).(T), nil

	default:
		return zero, fmt.Errorf("unsupported type %T for path value conversion", zero)
	}
}
