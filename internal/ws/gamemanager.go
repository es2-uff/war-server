package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"es2.uff/war-server/internal/domain/room"
	"github.com/google/uuid"
)

type GameManager struct {
	sync.RWMutex
	games map[string]*Game
}

type Gamelog struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Game struct {
	ID         string
	GameState  *GameState
	clients    map[*Client]bool
	log        []Gamelog
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
		log:        []Gamelog{},
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
					Armies:   0,
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
			log.Printf("Error processing next turn: %v", err)
		}
		playerName := g.GameState.Players[playerID].Username
		g.log = append(g.log, Gamelog{
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("%s finalizou o turno.", playerName),
		})
	case "attack":
		from, _ := msg["from"].(string)
		to, _ := msg["to"].(string)
		armiesFloat, _ := msg["attacking_armies"].(float64)
		if victory, err := g.GameState.Attack(playerID, from, to, int(armiesFloat)); err != nil {
			log.Printf("Error processing attack: %v", err)
		} else {
			playerName := g.GameState.Players[playerID].Username
			fromName := g.getTerritoryName(from)
			toName := g.getTerritoryName(to)
			logMessage := ""

			if victory {
				logMessage = "%s atacou de %s para %s com %d exércitos com sucesso."
			} else {
				logMessage = "%s atacou de %s para %s com %d exércitos e perdeu."
			}

			g.log = append(g.log, Gamelog{
				Timestamp: time.Now(),
				Message: fmt.Sprintf(
					logMessage,
					playerName,
					fromName,
					toName,
					int(armiesFloat),
				),
			})
		}
	case "troop_assign":
		territoryID, _ := msg["territory_id"].(string)
		if err := g.GameState.Deploy(playerID, territoryID); err != nil {
			log.Printf("Error processing deploy: %v", err)
		} else {
			playerName := g.GameState.Players[playerID].Username
			territoryName := g.getTerritoryName(territoryID)
			g.log = append(g.log, Gamelog{
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("%s posicionou 1 exército em %s.", playerName, territoryName),
			})
		}
	case "troop_move":
		from, _ := msg["from"].(string)
		to, _ := msg["to"].(string)
		armiesFloat, _ := msg["moving_armies"].(float64)
		if err := g.GameState.Move(playerID, from, to, int(armiesFloat)); err != nil {
			log.Printf("Error processing move: %v", err)
		} else {
			playerName := g.GameState.Players[playerID].Username
			fromName := g.getTerritoryName(from)
			toName := g.getTerritoryName(to)
			g.log = append(g.log, Gamelog{
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("%s moveu %d exércitos de %s para %s.", playerName, int(armiesFloat), fromName, toName),
			})
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
			"log":       g.log,
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

func (g *Game) getTerritoryName(territoryID string) string {
	g.GameState.RLock()
	defer g.GameState.RUnlock()

	for _, t := range g.GameState.Territories {
		if t.ID == territoryID {
			return t.Name
		}
	}
	return "Unknown"
}
