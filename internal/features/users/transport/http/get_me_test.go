package users_transport_http

import (
	"encoding/json"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMeHandler_Success(t *testing.T) {
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

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetMe(rr, req)

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

func TestGetMeHandler_NotFound(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		GetUser(mock.Anything, id).
		Return(domain.User{}, domain.NotFoundErr(domain.UserEntity, "id", id.String())).
		Once()

	handler := NewUsersHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetMe(rr, req)

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
			Message: domain.NotFoundErr(domain.UserEntity, "id", id.String()).Error(),
		},
	}

	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, want, response)
}
