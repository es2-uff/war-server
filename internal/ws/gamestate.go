package ws

import (
	"fmt"
	"slices"
	"sync"

	"es2.uff/war-server/internal/domain/battle"
	"es2.uff/war-server/internal/domain/card"
	"es2.uff/war-server/internal/domain/game"
	"es2.uff/war-server/internal/domain/objective"
	"es2.uff/war-server/internal/domain/player"
	"es2.uff/war-server/internal/domain/territory"
	"github.com/google/uuid"
)

type Player struct {
	ID            string       `json:"id"`
	Username      string       `json:"username"`
	Armies        int          `json:"armies"`
	Color         string       `json:"color"`
	IsReady       bool         `json:"is_ready"`
	IsOwner       bool         `json:"is_owner"`
	ObjectiveID   int          `json:"objective_id"`
	ObjectiveDesc string       `json:"objective_desc"`
	CardsInHand   []*card.Card `json:"cards_in_hand"`
}

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
	Deck                      *card.Deck         `json:"-"`
	TradesCount               int                `json:"trades_count"`
}

func NewGameState(roomID string) *GameState {
	gs := &GameState{
		RoomID:      roomID,
		Players:     make(map[string]*Player),
		Territories: nil,
		CurrentTurn: "",
		TradesCount: 2,
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

	territoryIDMap := make(map[int]string)
	gs.Territories = make([]*Territory, 0, len(domainTerritories))

	for _, dt := range domainTerritories {
		wsID := uuid.NewString()
		territoryIDMap[dt.TerritoryID] = wsID

		wsTerr := &Territory{
			ID:         wsID,
			Name:       territory.TerritoryNameMap[territory.TerritoryID(dt.TerritoryID)],
			Owner:      dt.OwnerID.String(),
			OwnerColor: dt.OwnerColor,
			Armies:     dt.ArmyQuantity,
			Adjacent:   []string{},
		}

		fmt.Println(*wsTerr)
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

	for _, domainPlayer := range domainPlayers {
		wsPlayer := gs.Players[domainPlayer.ID.String()]
		if wsPlayer != nil {
			wsPlayer.ObjectiveID = int(domainPlayer.ObjectiveID)
			if objDetails, exists := objective.ObjectiveDetails[domainPlayer.ObjectiveID]; exists {
				wsPlayer.ObjectiveDesc = objDetails.Description
			}
		}
	}

	gs.Deck = card.NewDeck()
	gs.Deck.Shuffle()

	gs.getTurnAdditionalTroopsLocked(domainPlayers[0].ID.String())
	gs.CurrentTurn = domainPlayers[0].ID.String()
}

func (gs *GameState) Move(playerID, fromTerritoryID, toTerritoryID string, movingArmies int) error {
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
		return fmt.Errorf("not the owner of moving territory")
	}

	if toTerritory.Owner != playerID {
		return fmt.Errorf("can only move troops to your own territory")
	}

	if fromTerritory.Armies <= movingArmies {
		return fmt.Errorf("not enough armies (must leave 1 for occupation)")
	}

	if movingArmies < 1 {
		return fmt.Errorf("must move at least 1 army")
	}

	if !slices.Contains(fromTerritory.Adjacent, toTerritoryID) {
		return fmt.Errorf("territories are not adjacent")
	}

	fromTerritory.Armies -= movingArmies
	toTerritory.Armies += movingArmies

	return nil
}

func (gs *GameState) Trade(playerID, card1, card2, card3 string) (int, error) {
	gs.Lock()
	defer gs.Unlock()

	player := gs.Players[playerID]
	if player == nil {
		return 0, fmt.Errorf("player not found")
	}

	cardNames := []string{card1, card2, card3}
	cardsToRemove := make([]*card.Card, 0, 3)

	for _, cardName := range cardNames {
		found := false
		for _, playerCard := range player.CardsInHand {
			if playerCard.TerritoryName == cardName {
				cardsToRemove = append(cardsToRemove, playerCard)
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf("card %s not found in player's hand", cardName)
		}
	}

	if len(cardsToRemove) != 3 {
		return 0, fmt.Errorf("must trade exactly 3 cards")
	}

	firstShape := cardsToRemove[0].Shape
	for _, c := range cardsToRemove {
		if c.Shape != firstShape {
			return 0, fmt.Errorf("all cards must have the same shape")
		}
	}

	newHand := make([]*card.Card, 0, len(player.CardsInHand)-3)
	for _, playerCard := range player.CardsInHand {
		shouldRemove := false
		for _, removeCard := range cardsToRemove {
			if playerCard.TerritoryName == removeCard.TerritoryName {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			newHand = append(newHand, playerCard)
		}
	}
	player.CardsInHand = newHand

	for _, c := range cardsToRemove {
		gs.Deck.AddToBottom(*c)
	}

	troopsReceived := 2 * gs.TradesCount
	player.Armies += troopsReceived
	gs.TradesCount++

	return troopsReceived, nil
}

func (gs *GameState) Attack(playerID, fromTerritoryID, toTerritoryID string, attackingArmies int) (bool, error) {
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
		return false, fmt.Errorf("territory not found")
	}

	if fromTerritory.Owner != playerID {
		return false, fmt.Errorf("not the owner of attacking territory")
	}

	if toTerritory.Owner == playerID {
		return false, fmt.Errorf("cannot attack your own territory")
	}

	if fromTerritory.Armies <= attackingArmies {
		return false, fmt.Errorf("not enough armies (must leave 1 for occupation)")
	}

	if attackingArmies > 3 || attackingArmies < 1 {
		return false, fmt.Errorf("attacking armies must be between 1 and 3")
	}

	if !slices.Contains(fromTerritory.Adjacent, toTerritoryID) {
		return false, fmt.Errorf("territories are not adjacent")
	}

	defendingArmies := min(toTerritory.Armies, 3)

	attackerDice := battle.RollDice(attackingArmies)
	defenderDice := battle.RollDice(defendingArmies)

	attackerLosses, defenderLosses := battle.CompareDice(attackerDice, defenderDice)
	attackResult := attackerLosses < defenderLosses

	fromTerritory.Armies -= attackerLosses
	toTerritory.Armies -= defenderLosses

	if toTerritory.Armies == 0 {
		toTerritory.Owner = playerID
		toTerritory.OwnerColor = fromTerritory.OwnerColor
		toTerritory.Armies = attackingArmies - attackerLosses
		fromTerritory.Armies -= (attackingArmies - attackerLosses)

		drawnCard := gs.Deck.Draw()
		if drawnCard != nil {
			gs.Players[playerID].CardsInHand = append(gs.Players[playerID].CardsInHand, drawnCard)
		}
	}

	return attackResult, nil
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
