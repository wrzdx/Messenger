package http_middleware

import (
	"fmt"
	core_errors "messenger/internal/core/errors"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := http_response.NewHTTPSender(log, w)
			defer func() {
				if p := recover(); p != nil {
					err := core_errors.Error{
						Err:     fmt.Errorf("during handle HTTP request got unexpected panic: %v", p),
						Code:    core_errors.INTERNAL_ERROR,
						Message: "Internal Server Error",
					}
					responseHandler.Error(err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
