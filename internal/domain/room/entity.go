package room

import (
	"github.com/google/uuid"
)

var OpenRooms []*Room

type Room struct {
	RoomID     uuid.UUID
	Name       string
	OwnerID    uuid.UUID
	OwnerName  string
	players    []uuid.UUID
	PlayerCount int
	MaxPlayers int
}

func NewRoom(name string, ownerID uuid.UUID, ownerName string) (*Room, error) {
	newRoom := &Room{
		RoomID:      uuid.New(),
		OwnerID:     ownerID,
		OwnerName:   ownerName,
		Name:        name,
		PlayerCount: 0,
		MaxPlayers:  6,
		players:     []uuid.UUID{},
	}

	OpenRooms = append(OpenRooms, newRoom)
	return newRoom, nil
}

func ListRooms() ([]*Room, error) {
	return OpenRooms, nil
}

func DeleteRoom(roomID uuid.UUID) {
	for i, room := range OpenRooms {
		if room.RoomID == roomID {
			OpenRooms = append(OpenRooms[:i], OpenRooms[i+1:]...)
			return
		}
	}
}

func GetRoom(roomID uuid.UUID) *Room {
	for _, room := range OpenRooms {
		if room.RoomID == roomID {
			return room
		}
	}
	return nil
}

func UpdatePlayerCount(roomID uuid.UUID, count int) {
	room := GetRoom(roomID)
	if room != nil {
		room.PlayerCount = count
	}
}
