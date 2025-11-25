package ws

import (
	"sync"
	"time"

	"es2.uff/war-server/internal/domain/game"
	"es2.uff/war-server/internal/domain/objective"
	"es2.uff/war-server/internal/domain/player"
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
	RoomID      string             `json:"room_id"`
	Players     map[string]*Player `json:"players"`
	Territories []*Territory       `json:"territories"`
	CurrentTurn string             `json:"current_turn"`
	Phase       string             `json:"phase"`
	OwnerID     string             `json:"owner_id"`
	LastUpdate  time.Time          `json:"last_update"`
}

func NewGameState(roomID string) *GameState {
	gs := &GameState{
		RoomID:      roomID,
		Players:     make(map[string]*Player),
		Territories: nil,
		Phase:       "waiting",
		LastUpdate:  time.Now(),
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

	gs.Territories = make([]*Territory, 0, len(domainTerritories))
	for _, dt := range domainTerritories {
		wsTerr := &Territory{
			ID:         uuid.NewString(),
			Name:       getTerritoryName(dt.TerritoryID),
			Owner:      dt.OwnerID.String(),
			OwnerColor: dt.OwnerColor,
			Armies:     dt.ArmyQuantity,
			Adjacent:   []string{}, // TODO: implement adjacency
		}
		gs.Territories = append(gs.Territories, wsTerr)
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

	gs.Phase = "deploy"
	gs.LastUpdate = time.Now()
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

func (gs *GameState) Attack(playerID, fromTerritoryID, toTerritoryID string, armies int) error {
	gs.Lock()
	defer gs.Unlock()

	// Find territories
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
		return nil // Territories not found
	}

	if fromTerritory.Owner != playerID {
		return nil // Not the owner
	}

	if fromTerritory.Armies <= armies {
		return nil // Not enough armies
	}

	// Simple battle: compare armies
	if armies > toTerritory.Armies {
		// Attacker wins
		toTerritory.Owner = playerID
		toTerritory.Armies = armies - toTerritory.Armies
		fromTerritory.Armies -= armies
	} else {
		// Defender wins
		fromTerritory.Armies -= armies
		toTerritory.Armies -= armies
	}

	gs.LastUpdate = time.Now()
	return nil
}

func (gs *GameState) Deploy(playerID, territoryID string, armies int) error {
	gs.Lock()
	defer gs.Unlock()

	player := gs.Players[playerID]
	if player == nil || player.Armies < armies {
		return nil // Not enough armies
	}

	// Find territory
	var territory *Territory
	for _, t := range gs.Territories {
		if t.ID == territoryID {
			territory = t
			break
		}
	}

	if territory == nil {
		return nil // Territory doesn't exist
	}

	// If territory is unclaimed, claim it
	if territory.Owner == "" || territory.Owner == playerID {
		territory.Owner = playerID
		territory.Armies += armies
		player.Armies -= armies
	}

	gs.LastUpdate = time.Now()
	return nil
}

func (gs *GameState) NextTurn() {
	gs.Lock()
	defer gs.Unlock()

	if len(gs.Players) == 0 {
		return
	}

	// Find next player
	found := false
	for pid := range gs.Players {
		if found {
			gs.CurrentTurn = pid
			gs.LastUpdate = time.Now()
			return
		}
		if pid == gs.CurrentTurn {
			found = true
		}
	}

	// If we didn't find a next player, go back to first
	for pid := range gs.Players {
		gs.CurrentTurn = pid
		break
	}

	gs.LastUpdate = time.Now()
}
