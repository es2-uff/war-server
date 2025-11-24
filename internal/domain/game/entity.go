package game

import (
	"math/rand/v2"

	"es2.uff/war-server/internal/domain/objective"
	"es2.uff/war-server/internal/domain/player"
	"es2.uff/war-server/internal/domain/territory"
	"github.com/google/uuid"
)

type Game struct {
	RoomID uuid.UUID

	PlayerTurn      uuid.UUID
	TerritoriesList []int
}

func InstantiateGameTerritories(players []*player.Player) []*territory.Territory {
	var tl []*territory.Territory

	for _, territoryID := range territory.AllTerritories {
		region := territory.TerritoryRegionMap[territoryID]
		newTerritory := &territory.Territory{
			TerritoryID:  int(territoryID),
			RegionID:     int(region),
			OwnerID:      uuid.Nil,
			ArmyQuantity: 1,
		}
		tl = append(tl, newTerritory)
	}

	playerIterator := 0
	copiedTerritoriesList := make([]*territory.Territory, len(tl))
	copy(copiedTerritoriesList, tl)

	for len(copiedTerritoriesList) > 0 {
		randomIndex := rand.IntN(len(copiedTerritoriesList))
		randomTerritory := copiedTerritoriesList[randomIndex]

		randomTerritory.OwnerID = players[playerIterator].ID
		randomTerritory.ArmyQuantity = 1

		playerIterator++
		if playerIterator >= len(players) {
			playerIterator = 0
		}

		copiedTerritoriesList = append(
			copiedTerritoriesList[:randomIndex],
			copiedTerritoriesList[randomIndex+1:]...,
		)
	}

	return tl
}

func AssignObjectivesToPlayers(players []*player.Player) {
	availableObjectives := make([]objective.ObjectiveID, len(objective.AllObjectives))
	copy(availableObjectives, objective.AllObjectives)

	for _, p := range players {
		randomIndex := rand.IntN(len(availableObjectives))
		randomObjective := availableObjectives[randomIndex]

		p.ObjectiveID = randomObjective

		availableObjectives = append(
			availableObjectives[:randomIndex],
			availableObjectives[randomIndex+1:]...,
		)
	}
}
