package users_transport_http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	core_context "messenger/internal/core/context"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMe(t *testing.T) {
	user := newUsersTransportTestUser(t)
	service := NewMockUsersService(t)
	service.EXPECT().
		GetUser(mock.Anything, user.ID).
		Return(user, nil)
	handler := NewUsersHandler(service)
	request := newUsersTransportRequest(t, http.MethodGet, "/users/me")
	ctx := core_context.WithClaims(
		request.Context(),
		core_context.ContextClaims{UserID: user.ID},
	)
	request = request.WithContext(ctx)
	recorder := httptest.NewRecorder()

	handler.GetMe(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, userDTOFromDomain(user), decodeUsersTransportData(t, recorder))
	require.NotContains(t, recorder.Body.String(), user.PasswordHash)
}
