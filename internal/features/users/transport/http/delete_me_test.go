package users_transport_http

import (
	"encoding/json"
	"errors"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteMeHandler_Success(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		DeleteUser(mock.Anything, id).
		Return(nil).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/users/me", nil)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.DeleteMe(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestDeleteMeHandler_NotFound(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		DeleteUser(mock.Anything, id).
		Return(domain.NotFoundErr(domain.UserEntity, "id", id.String())).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/users/me", nil)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.DeleteMe(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

}

func TestDeleteMeHandler_InternalError(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		DeleteUser(mock.Anything, id).
		Return(errors.New("database error")).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/users/me", nil)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.DeleteMe(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

}
