package auth_transport_http

import (
	"encoding/json"
	"messenger/internal/core/auth"
	core_errors "messenger/internal/core/errors"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRefreshHandler_Success(t *testing.T) {
	service := NewMockAuthService(t)

	expectedTokens := auth.TokenPair{
		Access:  "new-access",
		Refresh: "new-refresh",
	}

	cookie := NewMockCookieManager(t)

	cookie.EXPECT().
		GetRefreshToken(mock.Anything).
		Return("old-refresh", nil).
		Once()

	service.EXPECT().
		Refresh(mock.Anything, "old-refresh").
		Return(expectedTokens, nil).
		Once()

	cookie.EXPECT().
		SetRefreshToken(mock.Anything, "new-refresh").
		Once()

	handler := NewAuthHTTPHandler(service, cookie)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Refresh(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var response struct {
		Success bool            `json:"success"`
		Data    RefreshResponse `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "new-access", response.Data.Access)
}

func TestRefreshHandler_NoCookie(t *testing.T) {
	service := NewMockAuthService(t)

	cookie := NewMockCookieManager(t)

	cookie.EXPECT().
		GetRefreshToken(mock.Anything).
		Return("", auth.ErrInvalidToken).
		Once()

	handler := NewAuthHTTPHandler(service, cookie)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Refresh(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, core_errors.INVALID_TOKEN, response.Error.Code)
	assert.Equal(t, "get refresh token: invalid token", response.Error.Message)
}
