package users_transport_http

import (
	"encoding/json"
	"errors"
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

func TestGetUsersHandler_Success(t *testing.T) {
	service := NewMockUsersService(t)

	users := []domain.User{
		{
			Username:  "alice",
			FirstName: "Alice",
		},
		{
			Username:  "bob",
			FirstName: "Bob",
		},
	}

	service.EXPECT().
		GetUsers(mock.Anything, domain.Pagination{}).
		Return(users, nil).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Success bool             `json:"success"`
		Data    GetUsersResponse `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(
		t,
		GetUsersResponse(usersDTOFromDomains(users)),
		response.Data,
	)
}

func TestGetUsersHandler_WithPagination(t *testing.T) {
	service := NewMockUsersService(t)

	limit := 10
	offset := 20

	pagination := domain.NewPagination(&limit, &offset)

	service.EXPECT().
		GetUsers(mock.Anything, pagination).
		Return([]domain.User{}, nil).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users?limit=10&offset=20",
		nil,
	)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetUsersHandler_InvalidLimit(t *testing.T) {
	service := NewMockUsersService(t)

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users?limit=abc",
		nil,
	)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetUsersHandler_InvalidOffset(t *testing.T) {
	service := NewMockUsersService(t)

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users?offset=abc",
		nil,
	)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetUsersHandler_NegativeLimit(t *testing.T) {
	service := NewMockUsersService(t)
	service.EXPECT().GetUsers(mock.Anything, domain.NewPagination(new(-1), nil)).Return(nil, domain.ErrValidation)
	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users?limit=-1",
		nil,
	)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(
		t,
		domain.ErrValidation.Error(),
		response.Error.Message,
	)
}

func TestGetUsersHandler_NegativeOffset(t *testing.T) {
	service := NewMockUsersService(t)
	service.EXPECT().GetUsers(mock.Anything, domain.NewPagination(nil, new(-1))).Return(nil, domain.ErrValidation)
	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users?offset=-1",
		nil,
	)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(
		t,
		domain.ErrValidation.Error(),
		response.Error.Message,
	)
}

func TestGetUsersHandler_ServiceError(t *testing.T) {
	service := NewMockUsersService(t)

	service.EXPECT().
		GetUsers(mock.Anything, domain.Pagination{}).
		Return(nil, errors.New("database error")).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
