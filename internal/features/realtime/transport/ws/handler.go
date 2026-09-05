package realtime_transport_ws

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	authReqType = "authenticate"
	authResType = "authenticated"
)

type Handler struct {
	jwtProvider TokenProvider
	ctx         context.Context
}

func NewWSHandler(ctx context.Context, tp TokenProvider) *Handler {
	return &Handler{
		jwtProvider: tp,
		ctx:         ctx,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	var request AuthRequest
	authCtx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	err = wsjson.Read(authCtx, conn, &request)
	cancel()

	if err != nil {
		return
	}
	if request.Type != authReqType {
		conn.Close(websocket.StatusPolicyViolation, "authentication required")
		return
	}

	claims, err := h.jwtProvider.ParseAccessToken(request.AccessToken)
	if err != nil {
		conn.Close(websocket.StatusPolicyViolation, "invalid token")
		return
	}
	ctx, cancelExpired := context.WithDeadline(h.ctx, claims.ExpiresAt)
	defer cancelExpired()

	writeCtx, cancelWrite := context.WithTimeout(ctx, 5*time.Second)
	err = wsjson.Write(writeCtx, conn, AuthResponse{
		Type: authResType,
	})
	cancelWrite()

	if err != nil {
		return
	}

	closed := conn.CloseRead(ctx)
	<-closed.Done()
}

type AuthRequest struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	Type string `json:"type"`
}
