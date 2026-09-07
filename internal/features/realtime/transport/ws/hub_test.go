package realtime_transport_ws

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func queuedClient(t *testing.T, capacity int) *Client {
	t.Helper()
	client := NewClient(context.Background(), nil)
	client.send = make(chan []byte, capacity)
	t.Cleanup(client.Stop)
	return client
}

func TestHubPublishesToAllConnectionsOfOnlyTargetUser(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	first, second, other := queuedClient(t, 2), queuedClient(t, 2), queuedClient(t, 2)
	hub.Register(userID, first)
	hub.Register(userID, first) // Repeated registration must not duplicate delivery.
	hub.Register(userID, second)
	hub.Register(uuid.New(), other)
	event := Event{Type: "message.created", Data: map[string]string{"content": "hello"}}
	require.NoError(t, hub.Publish(userID, event))
	for _, client := range []*Client{first, second} {
		require.Len(t, client.send, 1)
		var got Event
		require.NoError(t, json.Unmarshal(<-client.send, &got))
		require.Equal(t, event.Type, got.Type)
		require.Equal(t, "hello", got.Data.(map[string]any)["content"])
	}
	require.Empty(t, other.send)
}

func TestHubUnregisterRemovesOnlyRequestedConnection(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	first, second := queuedClient(t, 2), queuedClient(t, 2)
	hub.Register(userID, first)
	hub.Register(userID, second)
	hub.Unregister(userID, first)
	hub.Unregister(userID, first)
	require.NoError(t, hub.Publish(userID, Event{Type: "test"}))
	require.Empty(t, first.send)
	require.Len(t, second.send, 1)
	hub.Unregister(userID, second)
	require.Empty(t, hub.clients)
	require.NoError(t, hub.Publish(userID, Event{Type: "test"}))
}

func TestHubStopsFullQueueWithoutBlockingHealthyConnection(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	slow, healthy := queuedClient(t, 1), queuedClient(t, 3)
	hub.Register(userID, slow)
	hub.Register(userID, healthy)
	require.NoError(t, hub.Publish(userID, Event{Type: "first"}))
	require.NoError(t, hub.Publish(userID, Event{Type: "second"}))
	require.ErrorIs(t, slow.ctx.Err(), context.Canceled)
	require.NoError(t, healthy.ctx.Err())
	require.Len(t, healthy.send, 2)
	// Even if an earlier publisher retained a snapshot, a stopped client rejects it.
	<-slow.send
	require.False(t, slow.enqueue([]byte(`{"type":"late"}`)))
	slow.Stop()
	slow.Stop()
}

func TestHubRejectsInvalidEventBeforeEnqueuing(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	client := queuedClient(t, 2)
	hub.Register(userID, client)
	require.Error(t, hub.Publish(userID, Event{}))
	require.Error(t, hub.Publish(userID, Event{Type: "test", Data: make(chan int)}))
	require.Empty(t, client.send)
}

func TestHubSerializesDataBeforePublishReturns(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	client := queuedClient(t, 1)
	hub.Register(userID, client)
	data := map[string]string{"content": "before"}
	require.NoError(t, hub.Publish(userID, Event{Type: "test", Data: data}))
	data["content"] = "after"
	require.JSONEq(t, `{"type":"test","data":{"content":"before"}}`, string(<-client.send))
}

func TestHubConcurrentPublishRegisterAndStop(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	var wg sync.WaitGroup
	errors := make(chan error, 4)
	for range 4 {
		wg.Go(func() {
			for range 100 {
				client := NewClient(context.Background(), nil)
				hub.Register(userID, client)
				if err := hub.Publish(userID, Event{Type: "test"}); err != nil {
					errors <- err
					return
				}
				client.Stop()
				hub.Unregister(userID, client)
			}
		})
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	require.Empty(t, hub.clients)
}
