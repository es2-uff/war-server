package ws

import (
	"log"
	"sync"

	"es2.uff/war-server/internal/domain/player"
	"es2.uff/war-server/internal/domain/room"
	"github.com/google/uuid"
)

type RoomServer struct {
	sync.RWMutex
	rooms map[string]*RoomHub
}

func NewRoomServer() *RoomServer {
	return &RoomServer{
		rooms: make(map[string]*RoomHub),
	}
}

func (rs *RoomServer) GetOrCreateHub(roomID string, userID string) *RoomHub {
	rs.Lock()
	defer rs.Unlock()

	if hub, exists := rs.rooms[roomID]; exists {
		playerUUID, _ := uuid.Parse(userID)
		roomUUID, _ := uuid.Parse(roomID)
		joiningPlayer := player.GetPlayer(playerUUID)

		room.AddPlayerToRoom(roomUUID, joiningPlayer)

		log.Printf("Player %s joined room %s", joiningPlayer.Name, roomID)
		return hub
	}

	hub := NewRoomHub(roomID, rs.handleOwnerLeft)
	rs.rooms[roomID] = hub

	go hub.Run()

	return hub
}

func (rs *RoomServer) GetHub(roomID string) *RoomHub {
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
