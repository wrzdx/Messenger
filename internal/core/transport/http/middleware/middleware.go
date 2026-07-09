package http_middleware

import "net/http"

type Middleware func(http.Handler) http.Handler
