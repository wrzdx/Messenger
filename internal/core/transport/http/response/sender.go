package http_response

import (
	"encoding/json"
	core_errors "messenger/internal/core/errors"
	logger "messenger/internal/core/logger"
	"net/http"

	"go.uber.org/zap"
)

type HTTPSender struct {
	log *logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPSender(log *logger.Logger, rw http.ResponseWriter) *HTTPSender {
	return &HTTPSender{
		log: log,
		rw:  rw,
	}
}

type APIResponse struct {
	Success bool            `json:"success"`
	Data    any             `json:"data,omitempty"`
	Error   *APIErrorDetail `json:"error,omitempty"`
}

type APIErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func MapErrorCodeToStatus(errorCode core_errors.ErrorCode) int {
	switch errorCode {
	case core_errors.VALIDATION_ERROR:
		return http.StatusBadRequest
	case core_errors.USER_ALREADY_EXISTS:
		return http.StatusConflict
	case core_errors.USER_NOT_FOUND:
		return http.StatusNotFound
	case core_errors.INVALID_CREDENTIALS,
		core_errors.INVALID_TOKEN:
		return http.StatusUnauthorized
	case core_errors.WRONG_PASSWORD:
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}

func (s *HTTPSender) Error(err core_errors.Error) {
	statusCode := MapErrorCodeToStatus(err.Code)
	s.json(false, statusCode, nil, err)
}

func (s *HTTPSender) OK(statusCode int, data any) {
	s.json(true, statusCode, data, core_errors.Error{})
}

func (s *HTTPSender) json(success bool, statusCode int, data any, err core_errors.Error) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(statusCode)

	var resp APIResponse
	resp.Success = success
	if success {
		resp.Data = data
	} else {
		resp.Error = &APIErrorDetail{
			Code:    string(err.Code),
			Message: err.Message,
			Fields:  err.Fields,
		}
	}

	if err := json.NewEncoder(s.rw).Encode(resp); err != nil {
		s.log.Error("write HTTP response", zap.Error(err))
	}
}
