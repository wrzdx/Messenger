package auth

import "net/http"

type CookieManager interface {
	SetRefreshToken(
		w http.ResponseWriter,
		token string,
	)

	ClearRefreshToken(
		w http.ResponseWriter,
	)

	GetRefreshToken(
		r *http.Request,
	) (string, error)
}
