package http_response

type HTTPError struct {
	StatusCode int
	Code       string
	Message    string
	Fields     map[string]string
}

type ErrorMapper func(error) HTTPError
