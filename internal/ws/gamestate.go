package ws

import (
	"fmt"
	"slices"
	"sync"

	"es2.uff/war-server/internal/domain/battle"
	"es2.uff/war-server/internal/domain/game"
	"es2.uff/war-server/internal/domain/objective"
	"es2.uff/war-server/internal/domain/player"
	"es2.uff/war-server/internal/domain/territory"
	"github.com/google/uuid"
)

type Player struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Armies        int    `json:"armies"`
	Color         string `json:"color"`
	IsReady       bool   `json:"is_ready"`
	IsOwner       bool   `json:"is_owner"`
	ObjectiveID   int    `json:"objective_id"`
	ObjectiveDesc string `json:"objective_desc"`
}

// Territory represents a territory on the game board
type Territory struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Owner      string   `json:"owner"` // Player ID
	OwnerColor string   `json:"owner_color"`
	Armies     int      `json:"armies"`
	Adjacent   []string `json:"adjacent"` // Adjacent territory IDs
}

type GameState struct {
	sync.RWMutex
	RoomID                    string             `json:"room_id"`
	Players                   map[string]*Player `json:"players"`
	FinishedInitialDeployment []string           `json:"finished_initial_deployment"`
	Territories               []*Territory       `json:"territories"`
	CurrentTurn               string             `json:"current_turn"` // Player ID whose turn it is
	OwnerID                   string             `json:"owner_id"`
}

func NewGameState(roomID string) *GameState {
	gs := &GameState{
		RoomID:      roomID,
		Players:     make(map[string]*Player),
		Territories: nil,
		CurrentTurn: "",
	}

	return gs
}

func (gs *GameState) StartGame() {
	gs.Lock()
	defer gs.Unlock()

	domainPlayers := make([]*player.Player, 0, len(gs.Players))
	for playerID := range gs.Players {
		playerUUID, err := uuid.Parse(playerID)
		if err != nil {
			continue
		}
		domainPlayer := &player.Player{
			ID:    playerUUID,
			Name:  gs.Players[playerID].Username,
			Color: gs.Players[playerID].Color,
		}
		domainPlayers = append(domainPlayers, domainPlayer)
	}

	domainTerritories := game.InstantiateGameTerritories(domainPlayers)

	// First pass: create territories with IDs
	territoryIDMap := make(map[int]string) // Maps domain TerritoryID to WebSocket UUID
	gs.Territories = make([]*Territory, 0, len(domainTerritories))

	for _, dt := range domainTerritories {
		wsID := uuid.NewString()
		territoryIDMap[dt.TerritoryID] = wsID

		wsTerr := &Territory{
			ID:         wsID,
			Name:       getTerritoryName(dt.TerritoryID),
			Owner:      dt.OwnerID.String(),
			OwnerColor: dt.OwnerColor,
			Armies:     dt.ArmyQuantity,
			Adjacent:   []string{}, // Will be populated in second pass
		}
		gs.Territories = append(gs.Territories, wsTerr)
	}

	// Second pass: populate adjacency using the TerritoryAdjacencyMap
	for i, dt := range domainTerritories {
		adjacentTerritoryIDs := territory.TerritoryAdjacencyMap[territory.TerritoryID(dt.TerritoryID)]
		adjacentWSIDs := make([]string, 0, len(adjacentTerritoryIDs))

		for _, adjID := range adjacentTerritoryIDs {
			if wsID, exists := territoryIDMap[int(adjID)]; exists {
				adjacentWSIDs = append(adjacentWSIDs, wsID)
			}
		}

		gs.Territories[i].Adjacent = adjacentWSIDs
	}

	game.AssignObjectivesToPlayers(domainPlayers)

	// Update WebSocket players with their objectives
	for _, domainPlayer := range domainPlayers {
		wsPlayer := gs.Players[domainPlayer.ID.String()]
		if wsPlayer != nil {
			wsPlayer.ObjectiveID = int(domainPlayer.ObjectiveID)
			if objDetails, exists := objective.ObjectiveDetails[domainPlayer.ObjectiveID]; exists {
				wsPlayer.ObjectiveDesc = objDetails.Description
			}
		}
	}

	gs.getTurnAdditionalTroopsLocked(domainPlayers[0].ID.String())
	gs.CurrentTurn = domainPlayers[0].ID.String()
}

func getTerritoryName(territoryID int) string {
	// Map territory IDs to names
	names := map[int]string{
		0: "Algeria", 1: "Egypt", 2: "Sudan", 3: "Congo", 4: "South Africa", 5: "Madagascar",
		6: "England", 7: "Iceland", 8: "Sweden", 9: "Moscow", 10: "Germany", 11: "Poland", 12: "Portugal",
		13: "Middle East", 14: "India", 15: "Vietnam", 16: "China", 17: "Aral", 18: "Omsk", 19: "Dudinka", 20: "Siberia", 21: "Tchita", 22: "Mongolia", 23: "Japan", 24: "Vladivostok",
		25: "Australia", 26: "New Guinea", 27: "Sumatra", 28: "Borneo",
		29: "Brazil", 30: "Argentina", 31: "Chile", 32: "Colombia",
		33: "Mexico", 34: "California", 35: "New York", 36: "Labrador", 37: "Ottawa", 38: "Vancouver", 39: "Mackenzie", 40: "Alaska", 41: "Greenland",
	}
	if name, ok := names[territoryID]; ok {
		return name
	}
	return "Unknown"
}

func (gs *GameState) Attack(playerID, fromTerritoryID, toTerritoryID string, attackingArmies int) error {
	gs.Lock()
	defer gs.Unlock()

	var fromTerritory, toTerritory *Territory
	for _, t := range gs.Territories {
		if t.ID == fromTerritoryID {
			fromTerritory = t
		}
		if t.ID == toTerritoryID {
			toTerritory = t
		}
	}

	if fromTerritory == nil || toTerritory == nil {
		return fmt.Errorf("territory not found")
	}

	if fromTerritory.Owner != playerID {
		return fmt.Errorf("not the owner of attacking territory")
	}

	if toTerritory.Owner == playerID {
		return fmt.Errorf("cannot attack your own territory")
	}

	if fromTerritory.Armies <= attackingArmies {
		return fmt.Errorf("not enough armies (must leave 1 for occupation)")
	}

	if attackingArmies > 3 || attackingArmies < 1 {
		return fmt.Errorf("attacking armies must be between 1 and 3")
	}

	if !slices.Contains(fromTerritory.Adjacent, toTerritoryID) {
		return fmt.Errorf("territories are not adjacent")
	}

	defendingArmies := min(toTerritory.Armies, 3)

	attackerDice := battle.RollDice(attackingArmies)
	defenderDice := battle.RollDice(defendingArmies)

	attackerLosses, defenderLosses := battle.CompareDice(attackerDice, defenderDice)

	fromTerritory.Armies -= attackerLosses
	toTerritory.Armies -= defenderLosses

	if toTerritory.Armies == 0 {
		toTerritory.Owner = playerID
		toTerritory.OwnerColor = fromTerritory.OwnerColor
		toTerritory.Armies = attackingArmies - attackerLosses
		fromTerritory.Armies -= (attackingArmies - attackerLosses)
	}

	return nil
}

func (gs *GameState) Deploy(playerID, territoryID string) error {
	gs.Lock()
	defer gs.Unlock()

	player := gs.Players[playerID]
	if player == nil || player.Armies == 0 {
		return nil
	}

	var territory *Territory
	for _, t := range gs.Territories {
		if t.ID == territoryID {
			territory = t
			break
		}
	}

	if territory == nil {
		return nil
	}

	if territory.Owner == playerID {
		territory.Armies += 1
		player.Armies -= 1
	}

	return nil
}

func (gs *GameState) NextTurn(senderID string) error {
	gs.Lock()
	defer gs.Unlock()

	if len(gs.Players) == 0 {
		return nil
	}

	if gs.CurrentTurn != senderID {
		return fmt.Errorf("Not the turn owner")
	}

	playerIDs := make([]string, 0, len(gs.Players))
	for pid := range gs.Players {
		playerIDs = append(playerIDs, pid)
	}

	currentIndex := -1
	for i, pid := range playerIDs {
		if pid == gs.CurrentTurn {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(playerIDs)
	nextPlayerID := playerIDs[nextIndex]

	gs.CurrentTurn = nextPlayerID
	gs.getTurnAdditionalTroopsLocked(nextPlayerID)
	return nil
}

func (gs *GameState) GetTurnAdditionalTroops(playerID string) {
	gs.Lock()
	defer gs.Unlock()
	gs.getTurnAdditionalTroopsLocked(playerID)
}

func (gs *GameState) getTurnAdditionalTroopsLocked(playerID string) {
	player := gs.Players[playerID]
	if player == nil {
		return
	}

	territoriesOwned := 0

	for _, t := range gs.Territories {
		if t.Owner == playerID {
			territoriesOwned++
		}
	}

	if territoriesOwned < 6 {
		player.Armies += 3
	} else {
		player.Armies += territoriesOwned / 2
	}
}
