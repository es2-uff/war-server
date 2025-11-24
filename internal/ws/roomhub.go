package ws

import (
	"encoding/json"
	"log"
)

type RoomHub struct {
	ID         string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewRoomHub(roomID string, onOwnerLeft func(string)) *RoomHub {
	return &RoomHub{
		ID:         roomID,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the room hub's main loop
func (h *RoomHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.broadcastRoomState()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.broadcastRoomState()
			}

		case message := <-h.broadcast:
			shouldBroadcastState := h.handleMessage(message)
			if shouldBroadcastState {
				h.broadcastRoomState()
			}
		}
	}
}

// handleMessage processes incoming messages for room operations
// Returns true if room state should be broadcasted after handling
func (h *RoomHub) handleMessage(message []byte) bool {
	var msg map[string]any
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return false
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return false
	}

	switch msgType {
	case "player_ready":
		playerID, _ := msg["player_id"].(string)
		ready, _ := msg["ready"].(bool)

		// Find the client and update their ready status
		for client := range h.clients {
			if client.id == playerID {
				client.ready = ready
				log.Printf("Player %s ready status set to %v in room %s", playerID, ready, h.ID)
				break
			}
		}
		return true

	case "start_game":
		log.Printf("Start game requested in room %s", h.ID)
		h.broadcastGameStart()
		return false
	}

	return false
}

func (h *RoomHub) broadcastRoomState() {
	playerList := make([]map[string]any, 0)
	for client := range h.clients {
		playerList = append(playerList, map[string]any{
			"id":    client.id,
			"name":  client.username,
			"ready": client.ready,
		})
	}

	message := map[string]any{
		"type":    "room_update",
		"room_id": h.ID,
		"players": playerList,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling room state: %v", err)
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

func (h *RoomHub) broadcastGameStart() {
	message := map[string]any{
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

func (h *RoomHub) GetRegisterChan() chan *Client {
	return h.register
}

func (h *RoomHub) GetUnregisterChan() chan *Client {
	return h.unregister
}

func (h *RoomHub) GetBroadcastChan() chan []byte {
	return h.broadcast
}
