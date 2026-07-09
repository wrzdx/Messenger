package auth_transport_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	auth_service "messenger/internal/features/auth/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	service := NewMockAuthService(t)
	mockUser := test_utils.MockUser
	tokens := auth.TokenPair{Access: "access", Refresh: "refresh"}
	service.EXPECT().
		Register(mock.Anything, auth_service.NewRegisterPayload(
			mockUser.Username,
			mockUser.FirstName,
			mockUser.LastName,
			mockUser.Bio,
			"fsociety",
		)).
		Return(mockUser, tokens, nil).
		Once()

	cookie := NewMockCookieManager(t)
	cookie.EXPECT().SetRefreshToken(mock.Anything, "refresh").Once()
	handler := NewAuthHTTPHandler(service, cookie)

	requestBody := map[string]string{
		"username":   mockUser.Username,
		"first_name": mockUser.FirstName,
		"last_name":  *mockUser.LastName,
		"bio":        *mockUser.Bio,
		"password":   "fsociety",
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(bodyBytes))
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var responseBody struct {
		Success bool             `json:"success"`
		Data    RegisterResponse `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, responseBody.Success, true)
	assert.Equal(t, tokens.Access, responseBody.Data.Access)
	assert.Equal(t, mockUser.ID, responseBody.Data.User.ID)
	assert.Equal(t, mockUser.FirstName, responseBody.Data.User.FirstName)
	assert.Equal(t, mockUser.LastName, responseBody.Data.User.LastName)
	assert.Equal(t, mockUser.Bio, responseBody.Data.User.Bio)
}

func TestRegister_Fail(t *testing.T) {
	service := NewMockAuthService(t)
	mockUser := test_utils.MockUser
	service.EXPECT().
		Register(mock.Anything, auth_service.NewRegisterPayload(
			mockUser.Username,
			mockUser.FirstName,
			mockUser.LastName,
			mockUser.Bio,
			"fsociety",
		)).
		Return(domain.User{}, auth.TokenPair{},
			domain.AlreadyExistsErr(
				domain.UserEntity,
				"username",
				mockUser.Username,
			),
		).
		Once()

	cookie := NewMockCookieManager(t)
	handler := NewAuthHTTPHandler(service, cookie)

	requestBody := map[string]string{
		"username":   mockUser.Username,
		"first_name": mockUser.FirstName,
		"last_name":  *mockUser.LastName,
		"bio":        *mockUser.Bio,
		"password":   "fsociety",
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(bodyBytes))
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)

	var responseBody struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, responseBody.Success, false)
	assert.Equal(t, core_errors.ALREADY_EXISTS, responseBody.Error.Code)
	assert.Equal(t,
		fmt.Sprintf(
			"user with username='%s' already exists",
			mockUser.Username,
		),
		responseBody.Error.Message,
	)
}

func TestRegister_Validation(t *testing.T) {
	service := NewMockAuthService(t)
	cookie := NewMockCookieManager(t)
	handler := NewAuthHTTPHandler(service, cookie)

	requestBody := map[string]string{
		"last_name": strings.Repeat("a", 65),
		"bio":       strings.Repeat("a", 71),
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(bodyBytes))
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.Register(rr, req)

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
			Code:    core_errors.VALIDATION_ERROR,
			Message: core_errors.ErrValidation.Error(),
			Fields: map[string]string{
				"username":   "username is required and username must be between 5 and 32 characters",
				"first_name": "first_name is required and first_name must be between 1 and 64 characters",
				"last_name":  "last_name cannot exceed 64 characters",
				"bio":        "bio cannot exceed 70 characters",
				"password":   "password is required and password must be between 8 and 32 characters",
			},
		},
	}

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, want, responseBody)
}
