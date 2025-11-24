package room

import (
	"fmt"

	"es2.uff/war-server/internal/domain/player"
	"github.com/google/uuid"
)

var OpenRooms []*Room

type Room struct {
	RoomID      uuid.UUID
	Name        string
	OwnerID     uuid.UUID
	OwnerName   string
	Players     []*player.Player
	PlayerCount int
	MaxPlayers  int
}

func NewRoom(name string, ownerID uuid.UUID, ownerName string) (*Room, error) {

	newRoom := &Room{
		RoomID:      uuid.New(),
		OwnerID:     ownerID,
		OwnerName:   ownerName,
		Name:        name,
		PlayerCount: 0,
		Players:     []*player.Player{},
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

func AddPlayerToRoom(roomID uuid.UUID, player *player.Player) {
	for _, room := range OpenRooms {
		if room.RoomID == roomID {
			room.Players = append(room.Players, player)
			break
		}
	}

	fmt.Println(OpenRooms)
}
