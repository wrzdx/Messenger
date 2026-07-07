package core_http_response

import (
	"encoding/json"
	"fmt"
	core_logger "messenger/internal/core/logger"
	"net/http"

	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(log *core_logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) HTMLResponse(html []byte) {
	h.rw.WriteHeader(http.StatusOK)
	h.rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err := h.rw.Write(html); err != nil {
		h.log.Error("write HTML HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, zap.Error(err))
	h.ErrorResponse(MapError(err))
}

func (h *HTTPResponseHandler) ErrorResponse(err Error) {
	var (
		logFunc func(string, ...zap.Field)
	)

	switch err.Status {
	case http.StatusBadRequest,
		http.StatusConflict,
		http.StatusForbidden,
		http.StatusUnauthorized:
		logFunc = h.log.Warn
	case http.StatusNotFound:
		logFunc = h.log.Debug
	default:
		logFunc = h.log.Error
	}

	logFunc(err.Error.Error())
	h.JSONResponse(ErrorResponse{Error: err.Message}, err.Status)
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}
