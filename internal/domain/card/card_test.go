package card

import (
	"testing"

	"es2.uff/war-server/internal/domain/territory"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	if deck == nil {
		t.Fatal("NewDeck() returned nil")
	}

	expectedSize := len(territory.AllTerritories)
	if deck.Size() != expectedSize {
		t.Errorf("NewDeck() size = %d, want %d", deck.Size(), expectedSize)
	}

	// Check all cards are unique territories
	seen := make(map[territory.TerritoryID]bool)
	for _, card := range deck.Cards {
		if seen[card.TerritoryID] {
			t.Errorf("Duplicate territory in deck: %v", card.TerritoryID)
		}
		seen[card.TerritoryID] = true

		// Check card has valid name and shape
		if card.TerritoryName == "" {
			t.Errorf("Card for territory %v has empty name", card.TerritoryID)
		}
		if card.Shape < 0 {
			t.Errorf("Card for territory %v has invalid shape", card.TerritoryID)
		}
	}
}

func TestDeck_Draw(t *testing.T) {
	deck := NewDeck()
	initialSize := deck.Size()

	card := deck.Draw()

	if card == nil {
		t.Fatal("Draw() returned nil when deck has cards")
	}

	if deck.Size() != initialSize-1 {
		t.Errorf("After Draw(), deck size = %d, want %d", deck.Size(), initialSize-1)
	}
}

func TestDeck_DrawEmptyDeck(t *testing.T) {
	deck := &Deck{Cards: []Card{}}

	card := deck.Draw()

	if card != nil {
		t.Errorf("Draw() on empty deck = %v, want nil", card)
	}
}

func TestDeck_AddToBottom(t *testing.T) {
	deck := NewDeck()
	initialSize := deck.Size()

	newCard := Card{
		TerritoryID:   territory.Brazil,
		TerritoryName: "Brasil",
		Shape:         territory.Square,
	}

	deck.AddToBottom(newCard)

	if deck.Size() != initialSize+1 {
		t.Errorf("After AddToBottom(), deck size = %d, want %d", deck.Size(), initialSize+1)
	}

	// Check card is at bottom
	lastCard := deck.Cards[len(deck.Cards)-1]
	if lastCard.TerritoryID != newCard.TerritoryID {
		t.Errorf("Card at bottom = %v, want %v", lastCard.TerritoryID, newCard.TerritoryID)
	}
}

func TestDeck_Shuffle(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()

	// Record original order of deck1
	originalOrder := make([]territory.TerritoryID, len(deck1.Cards))
	for i, card := range deck1.Cards {
		originalOrder[i] = card.TerritoryID
	}

	deck1.Shuffle()

	// Check size unchanged
	if deck1.Size() != deck2.Size() {
		t.Errorf("After Shuffle(), deck size changed from %d to %d", deck2.Size(), deck1.Size())
	}

	// Check order changed (very unlikely to be same after shuffle)
	sameOrder := true
	for i, card := range deck1.Cards {
		if card.TerritoryID != originalOrder[i] {
			sameOrder = false
			break
		}
	}

	if sameOrder && deck1.Size() > 1 {
		t.Log("Warning: Shuffle() resulted in same order (very unlikely but possible)")
	}

	// Check all cards still present
	seen := make(map[territory.TerritoryID]bool)
	for _, card := range deck1.Cards {
		seen[card.TerritoryID] = true
	}

	for _, origID := range originalOrder {
		if !seen[origID] {
			t.Errorf("After Shuffle(), territory %v is missing", origID)
		}
	}
}

func TestDeck_DrawAllCards(t *testing.T) {
	deck := NewDeck()
	initialSize := deck.Size()

	cardsDrawn := 0
	for deck.Size() > 0 {
		card := deck.Draw()
		if card == nil {
			t.Fatalf("Draw() returned nil when deck size = %d", deck.Size())
		}
		cardsDrawn++
	}

	if cardsDrawn != initialSize {
		t.Errorf("Drew %d cards, expected %d", cardsDrawn, initialSize)
	}

	// Try drawing from empty deck
	card := deck.Draw()
	if card != nil {
		t.Errorf("Draw() from empty deck = %v, want nil", card)
	}
}

func TestDeck_Size(t *testing.T) {
	tests := []struct {
		name        string
		setupDeck   func() *Deck
		expectedSize int
	}{
		{
			name:        "New deck",
			setupDeck:   NewDeck,
			expectedSize: len(territory.AllTerritories),
		},
		{
			name: "Empty deck",
			setupDeck: func() *Deck {
				return &Deck{Cards: []Card{}}
			},
			expectedSize: 0,
		},
		{
			name: "Deck after draw",
			setupDeck: func() *Deck {
				d := NewDeck()
				d.Draw()
				return d
			},
			expectedSize: len(territory.AllTerritories) - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := tt.setupDeck()
			if deck.Size() != tt.expectedSize {
				t.Errorf("Size() = %d, want %d", deck.Size(), tt.expectedSize)
			}
		})
	}
}
