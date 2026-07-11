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
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (s *HTTPSender) Error(err error) {
	errResp := core_errors.MapError(err)
	s.json(false, errResp.Code, nil, errResp)
}

func (s *HTTPSender) OK(statusCode int, data any) {
	s.json(true, statusCode, data, core_errors.Error{})
}

func (s *HTTPSender) json(success bool, statusCode int, data any, err core_errors.Error) {
	s.rw.Header().Set("Content-Type", "application/json")
	s.rw.WriteHeader(statusCode)
	if statusCode == http.StatusNoContent {
		return
	}

	var resp APIResponse
	resp.Success = success
	if success {
		resp.Data = data
	} else {
		if statusCode >= 500 {
			s.log.Error(
				"CRITICAL_SERVER_ERROR",
				zap.Error(err),
			)
		} else {
			s.log.Warn(err.Error())
		}
		resp.Error = &APIErrorDetail{
			Message: err.Message,
			Fields:  err.Details,
		}
	}

	if err := json.NewEncoder(s.rw).Encode(resp); err != nil {
		s.log.Error("write HTTP response", zap.Error(err))
	}
}
