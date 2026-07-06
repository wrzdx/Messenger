package core_http_response

type ErrorResponse struct {
	Error string `json:"error"   example:"Error description"`
}
