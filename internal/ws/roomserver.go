package ws

import (
	"sync"
)

type RoomServer struct {
	sync.RWMutex
	rooms map[string]*Hub
}

func NewRoomServer() *RoomServer {
	return &RoomServer{
		rooms: make(map[string]*Hub),
	}
}

// GetOrCreateHub returns an existing hub or creates a new one
func (rs *RoomServer) GetOrCreateHub(roomID string) *Hub {
	rs.Lock()
	defer rs.Unlock()

	// Check if hub already exists
	if hub, exists := rs.rooms[roomID]; exists {
		return hub
	}

	// Create new hub
	hub := NewHub(roomID)
	rs.rooms[roomID] = hub

	// Start the hub
	go hub.Run()

	return hub
}

// GetHub returns an existing hub or nil if not found
func (rs *RoomServer) GetHub(roomID string) *Hub {
	rs.RLock()
	defer rs.RUnlock()

	return rs.rooms[roomID]
}

// RemoveHub removes a hub from the server
func (rs *RoomServer) RemoveHub(roomID string) {
	rs.Lock()
	defer rs.Unlock()

	delete(rs.rooms, roomID)
}
