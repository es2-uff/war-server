package ws

import (
	"encoding/json"
	"log"
	"sync"

	"es2.uff/war-server/internal/domain/room"
	"github.com/google/uuid"
)

type GameManager struct {
	sync.RWMutex
	games map[string]*Game
}

type Game struct {
	ID         string
	GameState  *GameState
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*Game),
	}
}

func (gm *GameManager) GetOrCreateGame(roomID string) *Game {
	gm.Lock()
	defer gm.Unlock()

	if game, exists := gm.games[roomID]; exists {
		return game
	}

	game := &Game{
		ID:         roomID,
		GameState:  NewGameState(roomID),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	roomUUID, err := uuid.Parse(roomID)
	if err == nil {
		r := room.GetRoom(roomUUID)
		if r != nil && len(r.Players) > 0 {
			colors := []string{"#FF0000", "#0066FF", "#00CC00", "#FFD700", "#9933FF", "#FF6600"}

			// Add players to game state
			for i, p := range r.Players {
				game.GameState.Players[p.ID.String()] = &Player{
					ID:       p.ID.String(),
					Username: p.Name,
					Armies:   20,
					Color:    colors[i%len(colors)],
					IsReady:  true,
				}
			}

			game.GameState.StartGame()
		}
	}

	gm.games[roomID] = game
	go game.Run()

	return game
}

func (g *Game) Run() {
	for {
		select {
		case client := <-g.register:
			g.clients[client] = true
			g.broadcastGameState()

		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.send)

				g.broadcastGameState()
			}

		case message := <-g.broadcast:
			g.handleMessage(message)
			g.broadcastGameState()
		}
	}
}

func (g *Game) handleMessage(message []byte) {
	var msg map[string]any
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	playerID, _ := msg["player_id"].(string)

	switch msgType {
	case "finish_turn":
		if err := g.GameState.NextTurn(playerID); err != nil {
			log.Printf("Error processing attack: %v", err)
		}
	case "attack":
		from, _ := msg["from"].(string)
		to, _ := msg["to"].(string)
		armies, _ := msg["armies"].(float64)
		if err := g.GameState.Attack(playerID, from, to, int(armies)); err != nil {
			log.Printf("Error processing attack: %v", err)
		}

	case "deploy":
		territory, _ := msg["territory"].(string)
		armies, _ := msg["armies"].(float64)
		if err := g.GameState.Deploy(playerID, territory, int(armies)); err != nil {
			log.Printf("Error processing deploy: %v", err)
		}
	}
}

func (g *Game) broadcastGameState() {
	g.GameState.RLock()
	defer g.GameState.RUnlock()

	for client := range g.clients {
		player := g.GameState.Players[client.id]
		if player == nil {
			continue
		}

		message := map[string]any{
			"type":      "update",
			"gameState": g.GameState,
		}

		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling game state: %v", err)
			continue
		}

		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(g.clients, client)
		}
	}
}

func (g *Game) GetRegisterChan() chan *Client {
	return g.register
}

func (g *Game) GetUnregisterChan() chan *Client {
	return g.unregister
}

func (g *Game) GetBroadcastChan() chan []byte {
	return g.broadcast
}
