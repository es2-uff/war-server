package player

import (
	"es2.uff/war-server/internal/domain/objective"
	"github.com/google/uuid"
)

var Players []*Player

type Player struct {
	ID   uuid.UUID
	Name string

	ArmyLeft    int
	ObjectiveID objective.ObjectiveID
}

func NewPlayer(name string) (*Player, error) {
	newPlayer := &Player{
		ID:   uuid.New(),
		Name: name,
	}

	Players = append(Players, newPlayer)
	return newPlayer, nil
}

func GetPlayer(playerID uuid.UUID) *Player {
	for _, player := range Players {
		if player.ID == playerID {
			return player
		}
	}
	return nil
}
