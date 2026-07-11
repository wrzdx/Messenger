package auth_transport_http

import (
	"bytes"
	"encoding/json"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	service := NewMockAuthService(t)

	expectedTokens := auth.TokenPair{Access: "access", Refresh: "refresh"}
	service.EXPECT().
		Login(mock.Anything, "ecorp", "fsociety").
		Return(expectedTokens, nil).
		Once()

	cookie := NewMockCookieManager(t)
	cookie.EXPECT().SetRefreshToken(mock.Anything, "refresh").Once()
	handler := NewAuthHTTPHandler(service, cookie)

	requestBody := map[string]string{
		"username": "ecorp",
		"password": "fsociety",
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyBytes))
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var responseBody struct {
		Success bool          `json:"success"`
		Data    LoginResponse `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, responseBody.Success, true)
	assert.Equal(t, expectedTokens.Access, responseBody.Data.Access)
}

func TestLogin_Fail(t *testing.T) {
	service := NewMockAuthService(t)

	service.EXPECT().
		Login(mock.Anything, "ecorp", "fsociety").
		Return(auth.TokenPair{}, domain.ErrInvalidCredentials).
		Once()

	cookie := NewMockCookieManager(t)
	handler := NewAuthHTTPHandler(service, cookie)

	requestBody := map[string]string{
		"username": "ecorp",
		"password": "fsociety",
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyBytes))
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var responseBody struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	want := struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}{
		Success: false,
		Error: http_response.APIErrorDetail{
			Message: domain.ErrInvalidCredentials.Error(),
		},
	}

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, want, responseBody)
}

func TestLogin_Validation(t *testing.T) {
	service := NewMockAuthService(t)
	cookie := NewMockCookieManager(t)
	handler := NewAuthHTTPHandler(service, cookie)

	requestBody := map[string]string{}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyBytes))
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var responseBody struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	want := struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}{
		Success: false,
		Error: http_response.APIErrorDetail{
			Message: domain.ErrValidation.Error() + " request",
			Fields: map[string]string{
				"username": "username is required",
				"password": "password is required",
			},
		},
	}

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, want, responseBody)
}
