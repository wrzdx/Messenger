package realtime_transport_ws

import "messenger/internal/core/auth"

type TokenProvider interface {
	ParseAccessToken(tokenStr string) (auth.ParsedAccessToken, error)
}
