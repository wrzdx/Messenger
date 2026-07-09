package users_transport_http

import (
	"bytes"
	"encoding/json"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
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

func TestPatchMeHandler_Success(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	username := "fsociety"

	request := map[string]string{
		"username": username,
	}
	expectedPatch := domain.UserPatch{
		Username: domain.Nullable[string]{
			Value: &username,
			Set:   true,
		},
	}

	user := domain.User{
		ID:        id,
		Username:  username,
		FirstName: "Elliot",
		CreatedAt: time.Now().Round(0),
	}

	service.EXPECT().
		PatchUser(mock.Anything, id, expectedPatch).
		Return(user, nil).
		Once()

	handler := NewUsersHTTPHandler(service)

	body, err := json.Marshal(request)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/users/me",
		bytes.NewReader(body),
	)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.PatchMe(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Success bool              `json:"success"`
		Data    PatchUserResponse `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(
		t,
		PatchUserResponse(userDTOFromDomain(user)),
		response.Data,
	)
}

func TestPatchMeHandler_Validation(t *testing.T) {
	service := NewMockUsersService(t)

	handler := NewUsersHTTPHandler(service)

	username := "abc"

	request := map[string]string{
		"username": username,
	}
	body, err := json.Marshal(request)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/users/me",
		bytes.NewReader(body),
	)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: uuid.New(),
	})
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.PatchMe(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, core_errors.VALIDATION_ERROR, response.Error.Code)
	assert.Equal(
		t,
		domain.ErrInvalidUsername.Error(),
		response.Error.Fields["username"],
	)
}

func TestPatchMeHandler_ServiceError(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	username := "fsociety"

	request := map[string]string{
		"username": username,
	}
	expectedPatch := domain.UserPatch{
		Username: domain.Nullable[string]{
			Value: &username,
			Set:   true,
		},
	}
	service.EXPECT().
		PatchUser(
			mock.Anything,
			id,
			expectedPatch,
		).
		Return(
			domain.User{},
			domain.NotFoundErr(domain.UserEntity, "id", id.String()),
		).
		Once()

	handler := NewUsersHTTPHandler(service)

	body, err := json.Marshal(request)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/users/me",
		bytes.NewReader(body),
	)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.PatchMe(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	var response struct {
		Success bool                         `json:"success"`
		Error   http_response.APIErrorDetail `json:"error"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, core_errors.NOT_FOUND, response.Error.Code)
}
