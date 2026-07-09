package users_transport_http

// import (
// 	"encoding/json"
// 	"messenger/internal/core/domain"
// 	http_response "messenger/internal/core/transport/http/response"
// 	test_utils "messenger/internal/core/utils/test"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"strings"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// )

// func TestGetUsers(t *testing.T) {
// 	tests := []struct {
// 		name              string
// 		limit             string
// 		offset            string
// 		serviceUsers      []domain.User
// 		serviceErr        error
// 		wantServiceLimit  *int
// 		wantServiceOffset *int
// 		wantServiceCalled bool
// 		wantStatus        int
// 		wantError         string
// 	}{
// 		{
// 			name:              "return all users",
// 			wantServiceCalled: true,
// 			wantStatus:        http.StatusOK,
// 			serviceUsers:      test_utils.Users,
// 		},
// 		{
// 			name:              "limit users",
// 			limit:             "1",
// 			wantServiceLimit:  new(1),
// 			serviceUsers:      test_utils.Users[:1],
// 			wantServiceCalled: true,
// 			wantStatus:        http.StatusOK,
// 		},
// 		{
// 			name:              "negative limit",
// 			limit:             "-1",
// 			wantServiceCalled: true,
// 			wantServiceLimit:  new(-1),
// 			serviceErr:        domain.ErrNegativeLimit,
// 			wantError:         http_response.MapError(domain.ErrNegativeLimit).Message,
// 			wantStatus:        http.StatusBadRequest,
// 		},
// 		{
// 			name:              "offset users",
// 			offset:            "1",
// 			wantServiceOffset: new(1),
// 			serviceUsers:      test_utils.Users[1:],
// 			wantServiceCalled: true,
// 			wantStatus:        http.StatusOK,
// 		},
// 		{
// 			name:              "negative offset",
// 			offset:            "-1",
// 			wantServiceCalled: true,
// 			wantServiceOffset: new(-1),
// 			serviceErr:        domain.ErrNegativeOffset,
// 			wantError:         http_response.MapError(domain.ErrNegativeOffset).Message,
// 			wantStatus:        http.StatusBadRequest,
// 		},
// 		{
// 			name:              "limit offset users",
// 			limit:             "1",
// 			offset:            "1",
// 			wantServiceLimit:  new(1),
// 			wantServiceOffset: new(1),
// 			serviceUsers:      test_utils.Users[1:2],
// 			wantServiceCalled: true,

// 			wantStatus: http.StatusOK,
// 		},
// 		{
// 			name:              "empty users",
// 			limit:             "1",
// 			offset:            "2",
// 			wantServiceLimit:  new(1),
// 			wantServiceOffset: new(2),
// 			serviceUsers:      test_utils.Users[2:2],
// 			wantServiceCalled: true,
// 			wantStatus:        http.StatusOK,
// 		},
// 		{
// 			name:       "invalid limit",
// 			limit:      "asadf",
// 			wantError:  http_response.ErrInvalidArgument.Error(),
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:      "invalid offset",
// 			offset:    "asadf",
// 			wantError: http_response.ErrInvalidArgument.Error(),

// 			wantStatus: http.StatusBadRequest,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup
// 			wantUsers := make([]UserDTOResponse, len(tt.serviceUsers))
// 			for i, user := range tt.serviceUsers {
// 				wantUsers[i] = userDTOFromDomain(user)
// 			}
// 			var (
// 				serviceCalled    bool
// 				serviceGotLimit  *int
// 				serviceGotOffset *int
// 			)

// 			service := StubUsersService{
// 				GetUsersFn: func(limit, offset *int) ([]domain.User, error) {
// 					serviceCalled = true
// 					serviceGotLimit = limit
// 					serviceGotOffset = offset

// 					return tt.serviceUsers, tt.serviceErr
// 				},
// 			}

// 			handler := NewUsersHTTPHandler(&service)
// 			rec := httptest.NewRecorder()
// 			values := url.Values{}
// 			if tt.limit != "" {
// 				values.Set("limit", tt.limit)
// 			}
// 			if tt.offset != "" {
// 				values.Set("offset", tt.offset)
// 			}

// 			req := httptest.NewRequest(
// 				http.MethodGet,
// 				"/users?"+values.Encode(),
// 				nil,
// 			)

// 			ctx := test_utils.GetLoggerContext(req.Context())
// 			// action
// 			handler.GetUsers(rec, req.WithContext(ctx))

// 			// check
// 			if serviceCalled != tt.wantServiceCalled {
// 				t.Fatalf(
// 					"service called = %v, want %v",
// 					serviceCalled,
// 					tt.wantServiceCalled,
// 				)
// 			}

// 			if serviceCalled {
// 				if diff := cmp.Diff(tt.wantServiceLimit, serviceGotLimit); diff != "" {
// 					t.Fatalf("limit mismatch (-want +got):\n%s", diff)
// 				}
// 				if diff := cmp.Diff(tt.wantServiceOffset, serviceGotOffset); diff != "" {
// 					t.Fatalf("offset mismatch (-want +got):\n%s", diff)
// 				}
// 			}

// 			if rec.Code != tt.wantStatus {
// 				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatus)
// 			}

// 			if tt.wantError != "" {
// 				var gotError http_response.ErrorResponse
// 				if err := json.NewDecoder(rec.Body).Decode(&gotError); err != nil {
// 					t.Fatalf("unexpected error: %v", err)
// 				}

// 				if !strings.HasSuffix(gotError.Error, tt.wantError) {
// 					t.Fatalf(
// 						"ErrorResponse mismatch:\nwant: %s\ngot: %s",
// 						tt.wantError,
// 						gotError.Error,
// 					)
// 				}
// 			} else {
// 				var gotResponse []UserDTOResponse
// 				if err := json.NewDecoder(rec.Body).Decode(&gotResponse); err != nil {
// 					t.Fatalf("unexpected error: %v", err)
// 				}

// 				if diff := cmp.Diff(wantUsers, gotResponse); diff != "" {
// 					t.Fatalf("GetUsersResponse mismatch (-want +got):\n%s", diff)
// 				}
// 			}
// 		})
// 	}
// }
