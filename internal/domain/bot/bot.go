package bot

import (
	"github.com/google/uuid"
)

type Bot struct {
	ID    uuid.UUID
	Name  string
	Color string
	IsBot bool
}

func NewBot(name, color string) *Bot {
	return &Bot{
		ID:    uuid.New(),
		Name:  name,
		Color: color,
		IsBot: true,
	}
}
