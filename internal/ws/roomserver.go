package ws

import (
	"log"
	"sync"

	"es2.uff/war-server/internal/domain/room"
	"github.com/google/uuid"
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

func (rs *RoomServer) GetOrCreateHub(roomID string) *Hub {
	rs.Lock()
	defer rs.Unlock()

	if hub, exists := rs.rooms[roomID]; exists {
		return hub
	}

	hub := NewHub(roomID, rs.handleOwnerLeft)
	rs.rooms[roomID] = hub

	go hub.Run()

	return hub
}

func (rs *RoomServer) GetHub(roomID string) *Hub {
	rs.RLock()
	defer rs.RUnlock()

	return rs.rooms[roomID]
}

func (rs *RoomServer) RemoveHub(roomID string) {
	rs.Lock()
	defer rs.Unlock()

	delete(rs.rooms, roomID)
}

func (rs *RoomServer) handleOwnerLeft(roomID string) {
	log.Printf("Handling owner left for room %s", roomID)
	rs.RemoveHub(roomID)

	roomUUID, err := uuid.Parse(roomID)
	if err == nil {
		room.DeleteRoom(roomUUID)
	}
}
