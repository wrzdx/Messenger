package http_middleware

import (
	"fmt"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			sender := http_response.NewHTTPSender(log, w)
			defer func() {
				if p := recover(); p != nil {
					sender.Error(fmt.Errorf("during handle HTTP request got unexpected panic: %v", p))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
