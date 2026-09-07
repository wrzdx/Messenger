package realtime_transport_ws

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]map[*Client]struct{}),
	}
}

func (h *Hub) Register(userID uuid.UUID, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	connections := h.clients[userID]
	if connections == nil {
		connections = make(map[*Client]struct{})
		h.clients[userID] = connections
	}
	connections[client] = struct{}{}
}
func (h *Hub) Unregister(userID uuid.UUID, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if connections := h.clients[userID]; connections != nil {
		delete(connections, client)
		if len(connections) == 0 {
			delete(h.clients, userID)
		}
	}
}

// Publish queues an event for the user's current connections. Success means
// local dispatch only, not delivery or acknowledgement by the browser.
// Concurrent Publish calls have no guaranteed ordering relative to each other.
func (h *Hub) Publish(userID uuid.UUID, event Event) error {
	if event.Type == "" {
		return fmt.Errorf("publish event: missing type")
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	// Keep the registry lock short. Clients remain valid even if Unregister
	// removes them after this snapshot; enqueue checks whether they stopped.
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients[userID]))
	for connection := range h.clients[userID] {
		clients = append(clients, connection)
	}
	h.mu.RUnlock()
	for _, client := range clients {
		client.enqueue(payload)
	}
	return nil
}
