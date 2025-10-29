package ws

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Armies   int    `json:"armies"`
	Color    string `json:"color"`
	IsReady  bool   `json:"is_ready"`
	IsOwner  bool   `json:"is_owner"`
}

// Territory represents a territory on the game board
type Territory struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Owner    string   `json:"owner"` // Player ID
	Armies   int      `json:"armies"`
	Adjacent []string `json:"adjacent"` // Adjacent territory IDs
}

type GameState struct {
	sync.RWMutex
	RoomID      string                `json:"room_id"`
	Players     map[string]*Player    `json:"players"`
	Territories map[string]*Territory `json:"territories"`
	CurrentTurn string                `json:"current_turn"`
	Phase       string                `json:"phase"`
	OwnerID     string                `json:"owner_id"`
	LastUpdate  time.Time             `json:"last_update"`
	// Card deck for the game
	Deck    []Card `json:"deck"`
	Discard []Card `json:"discard"`
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

// Card represents a territory card in the deck
type Card struct {
	ID          string `json:"id"`
	TerritoryID string `json:"territory_id"`
	Type        string `json:"type"`  // infantry, cavalry, artillery
	Owner       string `json:"owner"` // player ID who holds the card
}

// GetStartingArmiesForPlayerCount returns the standard starting armies for Risk-like rules
func GetStartingArmiesForPlayerCount(n int) int {
	switch n {
	case 2:
		return 40
	case 3:
		return 35
	case 4:
		return 30
	case 5:
		return 25
	case 6:
		return 20
	default:
		return 20
	}
}

// InitializeDeck builds a simple deck using territories and cycling card types
func (gs *GameState) InitializeDeck() {
	gs.Lock()
	defer gs.Unlock()

	gs.Deck = make([]Card, 0, len(gs.Territories))
	types := []string{"triangle", "circle", "square"}
	i := 0
	for tid := range gs.Territories {
		c := Card{
			ID:          tid + "-card",
			TerritoryID: tid,
			Type:        types[i%len(types)],
			Owner:       "",
		}
		gs.Deck = append(gs.Deck, c)
		i++
	}
	gs.LastUpdate = time.Now()
}

// ShuffleDeck shuffles the draw deck
func (gs *GameState) ShuffleDeck() {
	gs.Lock()
	defer gs.Unlock()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(gs.Deck), func(i, j int) {
		gs.Deck[i], gs.Deck[j] = gs.Deck[j], gs.Deck[i]
	})
	gs.LastUpdate = time.Now()
}

// drawCard pops a card from the top of the deck; if deck empty, returns error
func (gs *GameState) drawCard() (*Card, error) {
	gs.Lock()
	defer gs.Unlock()

	if len(gs.Deck) == 0 {
		return nil, errors.New("deck is empty")
	}
	c := gs.Deck[0]
	gs.Deck = gs.Deck[1:]
	gs.LastUpdate = time.Now()
	return &c, nil
}

// DrawCard gives a card to a player, recording ownership
func (gs *GameState) DrawCard(playerID string) (*Card, error) {
	c, err := gs.drawCard()
	if err != nil {
		return nil, err
	}
	gs.Lock()
	defer gs.Unlock()
	c.Owner = playerID
	gs.Discard = append(gs.Discard, *c) // keep track of dealt cards
	gs.LastUpdate = time.Now()
	return c, nil
}

// DealInitialCards deals n cards to each player (n may be 0)
func (gs *GameState) DealInitialCards(n int) error {
	if n <= 0 {
		return nil
	}

	gs.RLock()
	players := make([]string, 0, len(gs.Players))
	for pid := range gs.Players {
		players = append(players, pid)
	}
	gs.RUnlock()

	for i := 0; i < n; i++ {
		for _, pid := range players {
			_, err := gs.DrawCard(pid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AssignTerritoriesRandomly assigns territories round-robin to players and places one army on each
func (gs *GameState) AssignTerritoriesRandomly() {
	gs.Lock()
	defer gs.Unlock()

	// collect territories and players
	tids := make([]string, 0, len(gs.Territories))
	for tid := range gs.Territories {
		tids = append(tids, tid)
	}
	pids := make([]string, 0, len(gs.Players))
	for pid := range gs.Players {
		pids = append(pids, pid)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tids), func(i, j int) { tids[i], tids[j] = tids[j], tids[i] })

	// distribute
	for i, tid := range tids {
		pid := pids[i%len(pids)]
		t := gs.Territories[tid]
		t.Owner = pid
		t.Armies = 1
		if player := gs.Players[pid]; player != nil {
			if player.Armies > 0 {
				player.Armies -= 1
			}
		}
	}

	gs.LastUpdate = time.Now()
}

// SetupMatch runs the full starting setup: sets starting armies, assigns territories, initializes/shuffles deck and deals initial cards
func (gs *GameState) SetupMatch(initialCardsPerPlayer int) error {
	gs.Lock()
	if len(gs.Players) < 3 || len(gs.Players) > 6 {
		gs.Unlock()
		return errors.New("player count must be between 3 and 6 for setup")
	}

	// set starting armies
	start := GetStartingArmiesForPlayerCount(len(gs.Players))
	for _, p := range gs.Players {
		p.Armies = start
	}
	gs.Unlock()

	// assign territories (gives 1 army per territory and subtracts from player's pool)
	gs.AssignTerritoriesRandomly()

	// initialize deck and deal initial cards
	gs.InitializeDeck()
	gs.ShuffleDeck()
	if err := gs.DealInitialCards(initialCardsPerPlayer); err != nil {
		return err
	}

	gs.Lock()
	gs.Phase = "deploy"
	// choose a random starting player to place troops first
	if len(gs.Players) > 0 {
		pids := make([]string, 0, len(gs.Players))
		for pid := range gs.Players {
			pids = append(pids, pid)
		}
		rand.Seed(time.Now().UnixNano())
		gs.CurrentTurn = pids[rand.Intn(len(pids))]
	}
	gs.LastUpdate = time.Now()
	gs.Unlock()

	return nil
}

func (gs *GameState) AddPlayer(playerID, username string, isOwner bool) {
	gs.Lock()
	defer gs.Unlock()

	colors := []string{"red", "blue", "green", "yellow", "purple", "orange"}
	color := colors[len(gs.Players)%len(colors)]

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: username,
		Armies:   20,
		Color:    color,
		IsReady:  false,
		IsOwner:  isOwner,
	}

	if isOwner {
		gs.OwnerID = playerID
	}

	if len(gs.Players) == 1 {
		gs.CurrentTurn = playerID
	}

	gs.LastUpdate = time.Now()
}

func (gs *GameState) SetPlayerReady(playerID string, ready bool) {
	gs.Lock()
	defer gs.Unlock()

	if player, exists := gs.Players[playerID]; exists {
		player.IsReady = ready
		gs.LastUpdate = time.Now()
	}
}

func (gs *GameState) AllPlayersReady() bool {
	gs.RLock()
	defer gs.RUnlock()

	if len(gs.Players) < 3 {
		return false
	}

	for _, player := range gs.Players {
		if !player.IsReady {
			return false
		}
	}
	return true
}

func (gs *GameState) StartGame() {
	gs.Lock()
	defer gs.Unlock()

	gs.Phase = "deploy"
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
