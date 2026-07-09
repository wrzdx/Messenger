package users_transport_http

import (
	"context"
	"encoding/json"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserHandler_Success(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()
	now := time.Now().Round(0)

	user := domain.User{
		ID:        id,
		Username:  "ecorp",
		FirstName: "Elliot",
		CreatedAt: now,
	}

	service.EXPECT().
		GetUser(mock.Anything, id).
		Return(user, nil).
		Once()

	handler := NewUsersHTTPHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Success bool            `json:"success"`
		Data    GetUserResponse `json:"data"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, GetUserResponse(userDTOFromDomain(user)), response.Data)
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	service := NewMockUsersService(t)

	handler := NewUsersHTTPHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, core_errors.VALIDATION_ERROR, response.Error.Code)
}

func TestGetUserHandler_NotFound(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		GetUser(mock.Anything, id).
		Return(domain.User{}, domain.NotFoundErr(domain.UserEntity, "id", id.String())).
		Once()

	handler := NewUsersHTTPHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	req = req.WithContext(test_utils.GetLoggerContext(req.Context()))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	want := struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}{
		Success: false,
		Error: http_response.APIErrorDetail{
			Code:    core_errors.NOT_FOUND,
			Message: domain.NotFoundErr(domain.UserEntity, "id", id.String()).Error(),
		},
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, want, response)
}
