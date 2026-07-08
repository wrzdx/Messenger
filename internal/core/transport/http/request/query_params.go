package http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
			return nil, fmt.Errorf("param='%s' by key='%s' not a valid integer: %w", param, key, err)
		}
		res := any(val).(T)
		return &res, nil

	case *time.Time:
		const layout = "2006-01-02"
		date, err := time.Parse(layout, param)
		if err != nil {
			return nil, fmt.Errorf("param='%s' by key='%s' not a valid date: %w", param, key, err)
		}
		res := any(date).(T)
		return &res, nil

	case *uuid.UUID:
		id, err := uuid.Parse(param)
		if err != nil {
			return nil, fmt.Errorf("param='%s' by key='%s' not a valid uuid: %w", param, key, err)
		}
		res := any(id).(T)
		return &res, nil

	default:
		return nil, fmt.Errorf("unsupported type %T for query param conversion", target)
	}
}
