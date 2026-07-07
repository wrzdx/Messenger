package auth_transport_http

import "net/http"

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.cookieManger.ClearRefreshToken(w)
}
