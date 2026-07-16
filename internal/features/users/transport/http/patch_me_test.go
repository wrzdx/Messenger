package users_transport_http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	core_context "messenger/internal/core/context"
	"messenger/internal/core/domain"
	"messenger/internal/core/logger"
	core_types "messenger/internal/core/types"
	users_service "messenger/internal/features/users/service"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPatchMe(t *testing.T) {
	t.Run("passes patch command to service and returns updated user", func(t *testing.T) {
		user := newUsersTransportTestUser(t)
		username := "Updated_user"
		firstName := "Updated name"
		bio := "Updated bio"
		expectedCommand := users_service.UpdateProfileCommand{
			Username:  &username,
			FirstName: &firstName,
			LastName: core_types.Nullable[string]{
				Set:   true,
				Value: nil,
			},
			Bio: core_types.Nullable[string]{
				Set:   true,
				Value: &bio,
			},
		}
		service := NewMockUsersService(t)
		service.EXPECT().
			UpdateProfile(mock.Anything, user.ID, expectedCommand).
			Return(user, nil)
		handler := NewUsersHandler(service)
		request := newPatchMeRequest(t, user, `{
			"username":"Updated_user",
			"first_name":"Updated name",
			"last_name":null,
			"bio":"Updated bio"
		}`)
		recorder := httptest.NewRecorder()

		handler.PatchMe(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, userDTOFromDomain(user), decodeUsersTransportData(t, recorder))
		require.NotContains(t, recorder.Body.String(), user.PasswordHash)
	})

	t.Run("accepts empty patch as no-op command", func(t *testing.T) {
		user := newUsersTransportTestUser(t)
		service := NewMockUsersService(t)
		service.EXPECT().
			UpdateProfile(mock.Anything, user.ID, users_service.UpdateProfileCommand{}).
			Return(user, nil)
		handler := NewUsersHandler(service)
		request := newPatchMeRequest(t, user, `{}`)
		recorder := httptest.NewRecorder()

		handler.PatchMe(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("rejects malformed request before calling service", func(t *testing.T) {
		user := newUsersTransportTestUser(t)
		service := NewMockUsersService(t)
		handler := NewUsersHandler(service)
		request := newPatchMeRequest(t, user, `{"bio":42}`)
		recorder := httptest.NewRecorder()

		handler.PatchMe(recorder, request)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.Equal(t, "invalid_request", decodeUsersTransportError(t, recorder).Code)
	})

	t.Run("maps service error", func(t *testing.T) {
		user := newUsersTransportTestUser(t)
		service := NewMockUsersService(t)
		service.EXPECT().
			UpdateProfile(mock.Anything, user.ID, users_service.UpdateProfileCommand{}).
			Return(domain.User{}, fmt.Errorf("update profile: %w", domain.ErrNotFound))
		handler := NewUsersHandler(service)
		request := newPatchMeRequest(t, user, `{}`)
		recorder := httptest.NewRecorder()

		handler.PatchMe(recorder, request)

		require.Equal(t, http.StatusNotFound, recorder.Code)
		require.Equal(t, "user_not_found", decodeUsersTransportError(t, recorder).Code)
	})
}

func newPatchMeRequest(t *testing.T, user domain.User, body string) *http.Request {
	t.Helper()

	request := httptest.NewRequest(http.MethodPatch, "/users/me", strings.NewReader(body))
	ctx := logger.WithLogger(request.Context(), logger.NewTestLogger())
	ctx = core_context.WithClaims(ctx, core_context.ContextClaims{UserID: user.ID})
	return request.WithContext(ctx)
}
