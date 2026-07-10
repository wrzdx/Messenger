package auth_transport_http

import "net/http"

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.cookieManger.ClearRefreshToken(w)
}
