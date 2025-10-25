package ws

import (
	"encoding/json"
	"fmt"
	"log"
)

type Hub struct {
	// Room ID
	ID string

	// Game state for this room
	GameState *GameState

	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// NewHub creates a new hub for a game room
func NewHub(roomID string) *Hub {
	return &Hub{
		ID:         roomID,
		GameState:  NewGameState(roomID),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("Client Registered: %s in room %s\n", client.id, h.ID)

			// Add player to game state
			h.GameState.AddPlayer(client.id, client.id)

			// Send current game state to new client
			h.sendGameState(client)

			// Broadcast updated state to all clients
			h.broadcastGameState()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Printf("Client Unregistered: %s from room %s\n", client.id, h.ID)

				// Remove player from game state
				h.GameState.RemovePlayer(client.id)

				// Broadcast updated state to remaining clients
				h.broadcastGameState()
			}

		case message := <-h.broadcast:
			// Process the message and update game state
			h.handleMessage(message)

			// Broadcast updated game state to all clients
			h.broadcastGameState()
		}
	}
}

// handleMessage processes incoming messages and updates game state
func (h *Hub) handleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "attack":
		playerID, _ := msg["player_id"].(string)
		from, _ := msg["from"].(string)
		to, _ := msg["to"].(string)
		armies, _ := msg["armies"].(float64)

		h.GameState.Attack(playerID, from, to, int(armies))

	case "deploy":
		playerID, _ := msg["player_id"].(string)
		territory, _ := msg["territory"].(string)
		armies, _ := msg["armies"].(float64)

		h.GameState.Deploy(playerID, territory, int(armies))

	case "next_turn":
		h.GameState.NextTurn()
	}
}

// broadcastGameState sends the current game state to all connected clients
func (h *Hub) broadcastGameState() {
	h.GameState.RLock()
	defer h.GameState.RUnlock()

	message := map[string]interface{}{
		"type":      "update",
		"gameState": h.GameState,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling game state: %v", err)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// sendGameState sends the current game state to a specific client
func (h *Hub) sendGameState(client *Client) {
	h.GameState.RLock()
	defer h.GameState.RUnlock()

	message := map[string]interface{}{
		"type":      "update",
		"gameState": h.GameState,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling game state: %v", err)
		return
	}

	select {
	case client.send <- data:
	default:
		close(client.send)
		delete(h.clients, client)
	}
}
