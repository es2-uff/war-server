package ws

import (
	"testing"

	"es2.uff/war-server/internal/domain/card"
	"es2.uff/war-server/internal/domain/territory"
)

func TestGameState_Deploy(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	// Setup player and territory
	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: "Test Player",
		Armies:   5,
		Color:    "#FF0000",
	}

	territoryID := "territory1"
	gs.Territories = []*Territory{
		{
			ID:     territoryID,
			Name:   "Test Territory",
			Owner:  playerID,
			Armies: 1,
		},
	}

	// Deploy to own territory
	err := gs.Deploy(playerID, territoryID)
	if err != nil {
		t.Errorf("Deploy() error = %v, want nil", err)
	}

	// Check armies decreased
	if gs.Players[playerID].Armies != 4 {
		t.Errorf("Player armies = %d, want 4", gs.Players[playerID].Armies)
	}

	// Check territory armies increased
	if gs.Territories[0].Armies != 2 {
		t.Errorf("Territory armies = %d, want 2", gs.Territories[0].Armies)
	}
}

func TestGameState_Deploy_NoArmies(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: "Test Player",
		Armies:   0,
		Color:    "#FF0000",
	}

	territoryID := "territory1"
	gs.Territories = []*Territory{
		{
			ID:     territoryID,
			Name:   "Test Territory",
			Owner:  playerID,
			Armies: 5,
		},
	}

	initialArmies := gs.Territories[0].Armies

	err := gs.Deploy(playerID, territoryID)
	if err != nil {
		t.Errorf("Deploy() with no armies should return nil, got %v", err)
	}

	// Territory armies should not change
	if gs.Territories[0].Armies != initialArmies {
		t.Errorf("Territory armies changed from %d to %d when player has no armies", initialArmies, gs.Territories[0].Armies)
	}
}

func TestGameState_Deploy_NotOwner(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"
	otherPlayerID := "player2"

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: "Test Player",
		Armies:   5,
		Color:    "#FF0000",
	}

	territoryID := "territory1"
	gs.Territories = []*Territory{
		{
			ID:     territoryID,
			Name:   "Test Territory",
			Owner:  otherPlayerID,
			Armies: 5,
		},
	}

	initialArmies := gs.Territories[0].Armies
	initialPlayerArmies := gs.Players[playerID].Armies

	err := gs.Deploy(playerID, territoryID)
	if err != nil {
		t.Errorf("Deploy() to non-owned territory should return nil, got %v", err)
	}

	// Nothing should change
	if gs.Territories[0].Armies != initialArmies {
		t.Errorf("Territory armies changed when deploying to non-owned territory")
	}
	if gs.Players[playerID].Armies != initialPlayerArmies {
		t.Errorf("Player armies changed when deploying to non-owned territory")
	}
}

func TestGameState_Move(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: "Test Player",
		Color:    "#FF0000",
	}

	fromID := "territory1"
	toID := "territory2"

	gs.Territories = []*Territory{
		{
			ID:       fromID,
			Name:     "From Territory",
			Owner:    playerID,
			Armies:   10,
			Adjacent: []string{toID},
		},
		{
			ID:       toID,
			Name:     "To Territory",
			Owner:    playerID,
			Armies:   5,
			Adjacent: []string{fromID},
		},
	}

	err := gs.Move(playerID, fromID, toID, 5)
	if err != nil {
		t.Errorf("Move() error = %v, want nil", err)
	}

	// Check armies moved
	if gs.Territories[0].Armies != 5 {
		t.Errorf("From territory armies = %d, want 5", gs.Territories[0].Armies)
	}
	if gs.Territories[1].Armies != 10 {
		t.Errorf("To territory armies = %d, want 10", gs.Territories[1].Armies)
	}
}

func TestGameState_Move_NotEnoughArmies(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: "Test Player",
		Color:    "#FF0000",
	}

	fromID := "territory1"
	toID := "territory2"

	gs.Territories = []*Territory{
		{
			ID:       fromID,
			Name:     "From Territory",
			Owner:    playerID,
			Armies:   5,
			Adjacent: []string{toID},
		},
		{
			ID:       toID,
			Name:     "To Territory",
			Owner:    playerID,
			Armies:   3,
			Adjacent: []string{fromID},
		},
	}

	// Try to move all armies (must leave 1)
	err := gs.Move(playerID, fromID, toID, 5)
	if err == nil {
		t.Error("Move() with all armies should return error, got nil")
	}
}

func TestGameState_Move_NotAdjacent(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	gs.Players[playerID] = &Player{
		ID:       playerID,
		Username: "Test Player",
		Color:    "#FF0000",
	}

	fromID := "territory1"
	toID := "territory2"

	gs.Territories = []*Territory{
		{
			ID:       fromID,
			Name:     "From Territory",
			Owner:    playerID,
			Armies:   10,
			Adjacent: []string{}, // Not adjacent
		},
		{
			ID:       toID,
			Name:     "To Territory",
			Owner:    playerID,
			Armies:   5,
			Adjacent: []string{},
		},
	}

	err := gs.Move(playerID, fromID, toID, 5)
	if err == nil {
		t.Error("Move() to non-adjacent territory should return error, got nil")
	}
}

func TestGameState_Trade(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	gs.Deck = card.NewDeck()
	gs.TradesCount = 2

	card1 := &card.Card{TerritoryName: "Brasil", Shape: territory.Square}
	card2 := &card.Card{TerritoryName: "Argentina", Shape: territory.Square}
	card3 := &card.Card{TerritoryName: "Chile", Shape: territory.Square}

	gs.Players[playerID] = &Player{
		ID:          playerID,
		Username:    "Test Player",
		Armies:      0,
		CardsInHand: []*card.Card{card1, card2, card3},
	}

	troopsReceived, err := gs.Trade(playerID, "Brasil", "Argentina", "Chile")
	if err != nil {
		t.Errorf("Trade() error = %v, want nil", err)
	}

	expectedTroops := 2 * 2 // 2 * TradesCount
	if troopsReceived != expectedTroops {
		t.Errorf("Trade() troops = %d, want %d", troopsReceived, expectedTroops)
	}

	// Check player received armies
	if gs.Players[playerID].Armies != expectedTroops {
		t.Errorf("Player armies = %d, want %d", gs.Players[playerID].Armies, expectedTroops)
	}

	// Check cards removed
	if len(gs.Players[playerID].CardsInHand) != 0 {
		t.Errorf("Player still has %d cards, want 0", len(gs.Players[playerID].CardsInHand))
	}

	// Check TradesCount incremented
	if gs.TradesCount != 3 {
		t.Errorf("TradesCount = %d, want 3", gs.TradesCount)
	}
}

func TestGameState_Trade_DifferentShapes(t *testing.T) {
	gs := NewGameState("test-room")
	playerID := "player1"

	gs.Deck = card.NewDeck()

	card1 := &card.Card{TerritoryName: "Brasil", Shape: territory.Square}
	card2 := &card.Card{TerritoryName: "Argentina", Shape: territory.Circle}
	card3 := &card.Card{TerritoryName: "Chile", Shape: territory.Triangle}

	gs.Players[playerID] = &Player{
		ID:          playerID,
		Username:    "Test Player",
		Armies:      0,
		CardsInHand: []*card.Card{card1, card2, card3},
	}

	_, err := gs.Trade(playerID, "Brasil", "Argentina", "Chile")
	if err == nil {
		t.Error("Trade() with different shapes should return error, got nil")
	}
}

func TestGameState_NextTurn(t *testing.T) {
	gs := NewGameState("test-room")
	player1ID := "player1"
	player2ID := "player2"

	gs.Players[player1ID] = &Player{
		ID:       player1ID,
		Username: "Player 1",
		Armies:   0,
	}
	gs.Players[player2ID] = &Player{
		ID:       player2ID,
		Username: "Player 2",
		Armies:   0,
		IsBot:    true,
	}

	gs.Territories = []*Territory{
		{ID: "t1", Owner: player1ID},
		{ID: "t2", Owner: player1ID},
		{ID: "t3", Owner: player1ID},
		{ID: "t4", Owner: player2ID},
		{ID: "t5", Owner: player2ID},
		{ID: "t6", Owner: player2ID},
	}

	gs.CurrentTurn = player1ID

	botID, err := gs.NextTurn(player1ID)
	if err != nil {
		t.Errorf("NextTurn() error = %v, want nil", err)
	}

	// Check turn changed
	if gs.CurrentTurn == player1ID {
		t.Error("CurrentTurn should have changed")
	}

	// Check bot ID returned if next player is bot
	if botID != player2ID {
		t.Errorf("NextTurn() botID = %s, want %s", botID, player2ID)
	}

	// Check player2 received armies
	if gs.Players[player2ID].Armies <= 0 {
		t.Errorf("Player 2 armies = %d, should have received armies", gs.Players[player2ID].Armies)
	}
}

func TestGameState_NextTurn_NotCurrentPlayer(t *testing.T) {
	gs := NewGameState("test-room")
	player1ID := "player1"
	player2ID := "player2"

	gs.Players[player1ID] = &Player{ID: player1ID, Username: "Player 1"}
	gs.Players[player2ID] = &Player{ID: player2ID, Username: "Player 2"}
	gs.CurrentTurn = player1ID

	_, err := gs.NextTurn(player2ID)
	if err == nil {
		t.Error("NextTurn() by non-current player should return error, got nil")
	}
}
