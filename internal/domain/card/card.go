package card

import (
	"math/rand"

	"es2.uff/war-server/internal/domain/territory"
)

type Card struct {
	TerritoryID   territory.TerritoryID
	TerritoryName string
	Shape         territory.Shape
}

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	deck := &Deck{
		Cards: make([]Card, 0, len(territory.AllTerritories)),
	}

	for _, territoryID := range territory.AllTerritories {
		card := Card{
			TerritoryID:   territoryID,
			TerritoryName: territory.TerritoryNameMap[territoryID],
			Shape:         territory.TerritoryShapeMap[territoryID],
		}
		deck.Cards = append(deck.Cards, card)
	}

	return deck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw() *Card {
	if len(d.Cards) == 0 {
		return nil
	}

	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card
}

func (d *Deck) AddToBottom(card Card) {
	d.Cards = append(d.Cards, card)
}

func (d *Deck) Size() int {
	return len(d.Cards)
}
