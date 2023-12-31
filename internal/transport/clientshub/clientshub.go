package clientshub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

// Hub is a very simple implementation of connections manager.
// It is a safe to use in concurrent code.
type Hub struct {
	clientsMutex sync.Mutex

	idToClient map[string]*Connection
}

// New creates a new *Hub.
func New() *Hub {
	return &Hub{
		clientsMutex: sync.Mutex{},
		idToClient:   map[string]*Connection{},
	}
}

// AddClient registers a new connection in the hub.
func (h *Hub) AddClient(clientID string, eventsCh chan<- json.RawMessage) *Connection {
	conn := h.newConnection(clientID, eventsCh)

	h.clientsMutex.Lock()
	h.idToClient[conn.clientID] = conn
	h.clientsMutex.Unlock()

	return conn
}

// HasClient checks if the client exists in the hub.
func (h *Hub) HasClient(clientID string) bool {
	h.clientsMutex.Lock()
	_, ok := h.idToClient[clientID]
	h.clientsMutex.Unlock()

	return ok
}

// RemoveClient unregisters the client and removes it from the hub.
func (h *Hub) RemoveClient(clientID string) bool {
	h.clientsMutex.Lock()
	defer h.clientsMutex.Unlock()

	if _, ok := h.idToClient[clientID]; !ok {
		return false
	}

	delete(h.idToClient, clientID)

	return true
}

// SendToClients sends the given payload to multiple clients.
func (h *Hub) SendToClients(ctx context.Context, payload any, clientIDs ...string) error {
	errSend := make([]error, 0, len(clientIDs))

	for _, clientID := range clientIDs {
		err := h.sendToClient(ctx, payload, clientID)
		if err != nil {
			errSend = append(errSend, err)
		}
	}

	return errors.Join(errSend...)
}

func (h *Hub) sendToClient(ctx context.Context, payload any, clientID string) error {
	h.clientsMutex.Lock()
	conn, ok := h.idToClient[clientID]
	h.clientsMutex.Unlock()

	if !ok {
		return semerr.NewNotFoundError(
			fmt.Errorf("%w: %s", semerr.Error("client is not found"), clientID),
		)
	}

	event, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case conn.events <- event:
		return nil
	}
}
