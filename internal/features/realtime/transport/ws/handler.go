package realtime_transport_ws

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Handler struct {
	jwtProvider TokenProvider
	ctx         context.Context
	hub         *Hub
}

func NewWSHandler(ctx context.Context, tp TokenProvider, hub *Hub) *Handler {
	return &Handler{
		jwtProvider: tp,
		ctx:         ctx,
		hub:         hub,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()
	conn.SetReadLimit(4096)

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

	client := NewClient(ctx, conn)
	// Register before acknowledging so events published after the client sees
	// authenticated cannot fall into a registration gap. Run starts after the
	// acknowledgement, so queued events cannot overtake it on the wire.
	h.hub.Register(claims.UserID, client)
	defer func() {
		client.Stop()
		h.hub.Unregister(claims.UserID, client)
	}()

	writeCtx, cancelWrite := context.WithTimeout(client.ctx, writeTimeout)
	err = wsjson.Write(writeCtx, conn, AuthResponse{
		Type: authResType,
	})
	cancelWrite()

	if err != nil {
		return
	}

	// Any exit ends this connection; deferred cleanup removes its registration.
	_ = client.Run()
}
