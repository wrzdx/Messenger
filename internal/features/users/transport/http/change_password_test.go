package users_transport_http

import (
	"bytes"
	"encoding/json"
	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestChangePasswordHandler_Success(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		ChangePassword(
			mock.Anything,
			id,
			"old-password",
			"new-password",
		).
		Return(nil).
		Once()

	handler := NewUsersHandler(service)

	body := map[string]string{
		"old_password": "old-password",
		"new_password": "new-password",
	}

	data, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/users/password",
		bytes.NewReader(data),
	)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.ChangePassword(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestChangePasswordHandler_Fail(t *testing.T) {
	service := NewMockUsersService(t)

	id := uuid.New()

	service.EXPECT().
		ChangePassword(
			mock.Anything,
			id,
			"old-password",
			"new-password",
		).
		Return(domain.ErrWrongPassword)

	handler := NewUsersHandler(service)

	body := map[string]string{
		"old_password": "old-password",
		"new_password": "new-password",
	}

	data, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/users/password",
		bytes.NewReader(data),
	)

	ctx := test_utils.GetLoggerContext(req.Context())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{
		UserID: id,
	})

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.ChangePassword(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}
