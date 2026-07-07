package auth_transport_http

import (
	"bytes"
	"encoding/json"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name              string
		serviceErr        error
		wantServiceCalled bool
		wantCookieCalled  bool
		wantStatus        int
		wantError         string
		body              LoginRequest
	}{
		{
			name:              "valid credentials",
			wantServiceCalled: true,
			wantCookieCalled:  true,
			wantStatus:        http.StatusCreated,
			body: LoginRequest{
				Username: "i.ivanov",
				Password: "password",
			},
		},
		{
			name:              "invalid credentials",
			wantServiceCalled: true,
			wantStatus:        http.StatusUnauthorized,
			wantError:         core_http_response.MapError(domain.ErrInvalidCredentials).Message,
			serviceErr:        domain.ErrInvalidCredentials,
			body: LoginRequest{
				Username: "i.ivanov",
				Password: "password",
			},
		},
		{
			name:       "missing username",
			wantStatus: http.StatusBadRequest,
			wantError:  "username is required",
			body: LoginRequest{
				Password: "password",
			},
		},
		{
			name:       "missing password",
			wantStatus: http.StatusBadRequest,
			wantError:  "password is required",
			body: LoginRequest{
				Username: "i.ivanov",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := LoginResponse{
				Access: "access-token",
			}

			var (
				serviceCalled bool
				gotUsername   string
				gotPassword   string

				cookieCalled bool
				gotRefresh   string
			)

			service := StubAuthService{
				LoginFn: func(
					username string,
					password string,
				) (core_auth.AuthTokens, error) {
					serviceCalled = true
					gotUsername = username
					gotPassword = password

					return core_auth.AuthTokens{
						Access:  "access-token",
						Refresh: "refresh-token",
					}, tt.serviceErr
				},
			}

			cookies := StubCookieManager{
				SetRefreshTokenFn: func(
					w http.ResponseWriter,
					refresh string,
				) {
					cookieCalled = true
					gotRefresh = refresh
				},
			}

			handler := NewAuthHTTPHandler(&service, &cookies)

			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewReader(body),
			)
			ctx := core_test_utils.GetLoggerContext(req.Context())

			handler.Login(rec, req.WithContext(ctx))

			if serviceCalled != tt.wantServiceCalled {
				t.Fatalf(
					"service called = %v, want %v",
					serviceCalled,
					tt.wantServiceCalled,
				)
			}

			if serviceCalled {
				if gotUsername != tt.body.Username {
					t.Fatalf(
						"username = %q, want %q",
						gotUsername,
						tt.body.Username,
					)
				}

				if gotPassword != tt.body.Password {
					t.Fatalf(
						"password = %q, want %q",
						gotPassword,
						tt.body.Password,
					)
				}
			}

			if cookieCalled != tt.wantCookieCalled {
				t.Fatalf(
					"cookie called = %v, want %v",
					cookieCalled,
					tt.wantCookieCalled,
				)
			}

			if cookieCalled && gotRefresh != "refresh-token" {
				t.Fatalf(
					"refresh token = %q, want %q",
					gotRefresh,
					"refresh-token",
				)
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d",
					rec.Code,
					tt.wantStatus,
				)
			}

			if tt.wantError != "" {
				var gotError core_http_response.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotError); err != nil {
					t.Fatal(err)
				}

				if gotError.Error != tt.wantError {
					t.Fatalf(
						"want error %q, got %q",
						tt.wantError,
						gotError.Error,
					)
				}
			} else {
				var got LoginResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(want, got); diff != "" {
					t.Fatalf("LoginResponse mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
