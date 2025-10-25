package ws

import (
	"sync"
	"time"
)

// Player represents a connected player in the game
type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Armies   int    `json:"armies"`
	Color    string `json:"color"`
}

// Territory represents a territory on the game board
type Territory struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Owner    string `json:"owner"` // Player ID
	Armies   int    `json:"armies"`
	Adjacent []string `json:"adjacent"` // Adjacent territory IDs
}

// GameState represents the shared game state
type GameState struct {
	sync.RWMutex
	RoomID      string                `json:"room_id"`
	Players     map[string]*Player    `json:"players"`
	Territories map[string]*Territory `json:"territories"`
	CurrentTurn string                `json:"current_turn"` // Player ID whose turn it is
	Phase       string                `json:"phase"`        // "deploy", "attack", "fortify"
	LastUpdate  time.Time             `json:"last_update"`
}

// NewGameState creates a new game state with initial territories
func NewGameState(roomID string) *GameState {
	gs := &GameState{
		RoomID:      roomID,
		Players:     make(map[string]*Player),
		Territories: make(map[string]*Territory),
		Phase:       "waiting",
		LastUpdate:  time.Now(),
	}

	// Initialize some sample territories
	gs.Territories["north"] = &Territory{
		ID:       "north",
		Name:     "Northern Lands",
		Owner:    "",
		Armies:   0,
		Adjacent: []string{"center", "east"},
	}
	gs.Territories["center"] = &Territory{
		ID:       "center",
		Name:     "Central Plains",
		Owner:    "",
		Armies:   0,
		Adjacent: []string{"north", "south", "east", "west"},
	}
	gs.Territories["south"] = &Territory{
		ID:       "south",
		Name:     "Southern Regions",
		Owner:    "",
		Armies:   0,
		Adjacent: []string{"center", "west"},
	}
	gs.Territories["east"] = &Territory{
		ID:       "east",
		Name:     "Eastern Territories",
		Owner:    "",
		Armies:   0,
		Adjacent: []string{"north", "center"},
	}
	gs.Territories["west"] = &Territory{
		ID:       "west",
		Name:     "Western Frontier",
		Owner:    "",
		Armies:   0,
		Adjacent: []string{"center", "south"},
	}

	return gs
}

// AddPlayer adds a player to the game
func (gs *GameState) AddPlayer(playerID, username string) {
	gs.Lock()
	defer gs.Unlock()

	colors := []string{"red", "blue", "green", "yellow", "purple", "orange"}
	color := colors[len(gs.Players)%len(colors)]

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: username,
		Armies:   20, // Starting armies
		Color:    color,
	}

	// If this is the first player, set them as current turn
	if len(gs.Players) == 1 {
		gs.CurrentTurn = playerID
		gs.Phase = "deploy"
	}

	gs.LastUpdate = time.Now()
}

// RemovePlayer removes a player from the game
func (gs *GameState) RemovePlayer(playerID string) {
	gs.Lock()
	defer gs.Unlock()

	delete(gs.Players, playerID)
	gs.LastUpdate = time.Now()

	// If the current turn player left, advance to next player
	if gs.CurrentTurn == playerID && len(gs.Players) > 0 {
		for pid := range gs.Players {
			gs.CurrentTurn = pid
			break
		}
	}
}

// Attack processes an attack action
func (gs *GameState) Attack(playerID, fromTerritoryID, toTerritoryID string, armies int) error {
	gs.Lock()
	defer gs.Unlock()

	// Simple attack logic - attacker wins if they have more armies
	fromTerritory := gs.Territories[fromTerritoryID]
	toTerritory := gs.Territories[toTerritoryID]

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

// Deploy adds armies to a territory
func (gs *GameState) Deploy(playerID, territoryID string, armies int) error {
	gs.Lock()
	defer gs.Unlock()

	player := gs.Players[playerID]
	if player == nil || player.Armies < armies {
		return nil // Not enough armies
	}

	territory := gs.Territories[territoryID]
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

// NextTurn advances to the next player's turn
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
