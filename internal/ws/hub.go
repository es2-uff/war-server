package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"es2.uff/war-server/internal/domain/room"
	"github.com/google/uuid"
)

type Hub struct {
	ID         string
	GameState  *GameState
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	onOwnerLeft func(roomID string)
}

func NewHub(roomID string, onOwnerLeft func(string)) *Hub {
	return &Hub{
		ID:          roomID,
		GameState:   NewGameState(roomID),
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		onOwnerLeft: onOwnerLeft,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("Client Registered: %s (%s) in room %s\n", client.username, client.id, h.ID)

			isOwner := client.isOwner
			h.GameState.AddPlayer(client.id, client.username, isOwner)

			if roomUUID, err := uuid.Parse(h.ID); err == nil {
				room.UpdatePlayerCount(roomUUID, len(h.GameState.Players))
			}

			h.sendGameState(client)
			h.broadcastGameState()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Printf("Client Unregistered: %s from room %s\n", client.id, h.ID)

				h.GameState.RemovePlayer(client.id)

				if client.isOwner && h.GameState.Phase == "waiting" {
					fmt.Printf("Owner left lobby %s before game started, closing lobby\n", h.ID)
					h.broadcastLobbyClosing()
					if h.onOwnerLeft != nil {
						h.onOwnerLeft(h.ID)
					}
					return
				}

				if roomUUID, err := uuid.Parse(h.ID); err == nil {
					room.UpdatePlayerCount(roomUUID, len(h.GameState.Players))
				}

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
	case "player_ready":
		playerID, _ := msg["player_id"].(string)
		ready, _ := msg["ready"].(bool)
		h.GameState.SetPlayerReady(playerID, ready)

	case "start_game":
		playerID, _ := msg["player_id"].(string)
		if playerID == h.GameState.OwnerID && h.GameState.AllPlayersReady() {
			h.GameState.StartGame()
			h.broadcastGameStart()
			return
		}

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

func (h *Hub) broadcastLobbyClosing() {
	message := map[string]interface{}{
		"type":    "lobby_closed",
		"message": "The lobby owner has left. Returning to lobby list.",
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling lobby closing message: %v", err)
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

func (h *Hub) broadcastGameStart() {
	message := map[string]interface{}{
		"type":    "game_started",
		"room_id": h.ID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling game start message: %v", err)
		return
	}

	log.Printf("Broadcasting game start for room %s", h.ID)

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
