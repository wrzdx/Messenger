package auth_transport_http

import (
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout(t *testing.T) {
	service := NewMockAuthService(t)
	cookie := NewMockCookieManager(t)
	cookie.EXPECT().ClearRefreshToken(mock.Anything).Once()
	handler := NewAuthHTTPHandler(service, cookie)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}
