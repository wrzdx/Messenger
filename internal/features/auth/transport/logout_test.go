package auth_transport_http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogout(t *testing.T) {
	var called bool

	cookies := StubCookieManager{
		ClearRefreshTokenFn: func(w http.ResponseWriter) {
			called = true
		},
	}

	handler := NewAuthHTTPHandler(&StubAuthService{}, &cookies)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	handler.Logout(rec, req)

	if !called {
		t.Fatal("ClearRefreshToken was not called")
	}
}
