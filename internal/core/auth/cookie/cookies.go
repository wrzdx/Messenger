package core_auth_cookie

import (
	"net/http"
	"time"
)

const refreshCookieName = "refresh_token"

type CookieManager struct {
	refreshTTL time.Duration
	secure     bool
	path       string
}

func NewCookieManager(
	refreshTTL time.Duration,
	secure bool,
	path string,
) *CookieManager {
	return &CookieManager{
		refreshTTL: refreshTTL,
		secure:     secure,
		path:       path,
	}
}

func (c *CookieManager) SetRefreshToken(
	w http.ResponseWriter,
	token string,
) {
	cookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     c.path,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(c.refreshTTL.Seconds()),
	}

	http.SetCookie(w, cookie)
}

func (c *CookieManager) ClearRefreshToken(
	w http.ResponseWriter,
) {
	cookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     c.path,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
}
