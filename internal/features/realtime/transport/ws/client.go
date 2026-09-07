package realtime_transport_ws

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	clientQueueSize = 64
	writeTimeout    = 5 * time.Second
)

// Client owns one connection. The handler calls Run once; the hub may enqueue
// concurrently. Stop is safe to call multiple times, including before Run.
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex // Serializes enqueue with Stop; never held during network I/O.
}

func NewClient(ctx context.Context, conn *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(ctx)
	return &Client{conn: conn, send: make(chan []byte, clientQueueSize), ctx: ctx, cancel: cancel}
}

// enqueue accepts an immutable JSON payload. A full queue disconnects a slow
// consumer instead of holding up publishers or growing memory without bounds.
func (c *Client) enqueue(payload []byte) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctx.Err() != nil {
		return false
	}
	select {
	case c.send <- payload:
		return true
	default:
		c.cancel()
		return false
	}
}

func (c *Client) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cancel()
	// Do not close send: a publisher may still hold a snapshot containing c.
}

func (c *Client) Run() error {
	defer c.Stop()
	ctx, cancel := context.WithCancel(c.ctx)
	readerDone := make(chan struct{})
	// Read processes Ping/Pong while waiting. Any application message after
	// authentication is rejected: commands still use HTTP.
	go func() {
		defer close(readerDone)
		defer cancel()
		if _, _, err := c.conn.Read(ctx); err == nil {
			_ = c.conn.Close(websocket.StatusPolicyViolation, "unexpected data message")
		}
	}()
	defer func() {
		cancel()
		_ = c.conn.CloseNow()
		<-readerDone
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case payload := <-c.send:
			if err := ctx.Err(); err != nil {
				return err
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := c.conn.Write(writeCtx, websocket.MessageText, payload)
			cancel()
			if err != nil {
				return fmt.Errorf("write event: %w", err)
			}
		}
	}
}
