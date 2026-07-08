package auth_cookie

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestClearRefreshToken(t *testing.T) {
	manager := NewCookieManager(time.Minute*60, false, "/api/v1/auth/refresh")

	rec := httptest.NewRecorder()

	manager.ClearRefreshToken(rec)

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]

	if cookie.Name != "refresh_token" {
		t.Fatalf("wrong token key")
	}

	if cookie.MaxAge != -1 {
		t.Fatalf("not cleared")
	}

	if cookie.Path != "/api/v1/auth/refresh" {
		t.Fatalf("correct cookie path")
	}
}
