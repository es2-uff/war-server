package room

import (
	"github.com/google/uuid"
)

var OpenRooms []*Room // TODO: This should not be a global var but we can fix this later

type Room struct {
	RoomID  uuid.UUID
	Name    string
	OwnerID uuid.UUID
	players []uuid.UUID
}

func NewRoom(name string) (*Room, error) {

	newRoom := &Room{
		RoomID:  uuid.New(),
		OwnerID: uuid.New(),
		Name:    name,
	}

	OpenRooms = append(OpenRooms, newRoom)
	return newRoom, nil
}

func ListRooms() ([]*Room, error) {
	return OpenRooms, nil
}
