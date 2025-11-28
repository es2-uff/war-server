package handlers

import (
	"log"
	"net/http"

	"es2.uff/war-server/internal/domain/player"
	"es2.uff/war-server/internal/domain/room"
	"es2.uff/war-server/internal/ws"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreatePlayerRequest struct {
	PlayerName string `json:"player_name"`
}

type CreatePlayerResponse struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

type CreateRoomRequest struct {
	RoomName string    `json:"room_name"`
	OwnerID  uuid.UUID `json:"owner_id"`
}

type JoinRoomRequest struct {
	RoomID   uuid.UUID `json:"room_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

type RoomResponse struct {
	RoomID      string `json:"room_id"`
	RoomName    string `json:"room_name"`
	OwnerID     string `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
	PlayerCount int    `json:"player_count"`
	MaxPlayers  int    `json:"max_players"`
}

type RoomHandler struct {
	roomServer *ws.RoomServer
}

func NewRoomHandler(roomServer *ws.RoomServer) *RoomHandler {
	return &RoomHandler{
		roomServer: roomServer,
	}
}

func (rh *RoomHandler) CreatePlayer(c echo.Context) error {
	r := new(CreatePlayerRequest)

	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON format")
	}

	newPlayer, err := player.NewPlayer(r.PlayerName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	response := CreatePlayerResponse{
		PlayerID:   newPlayer.ID.String(),
		PlayerName: newPlayer.Name,
	}

	log.Printf("New player %s created with ID %s", newPlayer.Name, newPlayer.ID)
	return c.JSON(http.StatusOK, response)
}

func (rh *RoomHandler) CreateNewRoom(c echo.Context) error {
	r := new(CreateRoomRequest)

	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON format")
	}

	owner := player.GetPlayer(r.OwnerID)

	nr, err := room.NewRoom(r.RoomName, owner.ID, owner.Name)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	// Add the owner to the room's player list
	nr.Players = append(nr.Players, owner)
	nr.PlayerCount = 1

	response := RoomResponse{
		RoomID: nr.RoomID.String(),
	}

	log.Printf("New room %s created successfully by %s (owner added to players)", nr.RoomID, nr.OwnerName)
	return c.JSON(http.StatusOK, response)
}

func (rh *RoomHandler) ListRooms(c echo.Context) error {
	resp := []RoomResponse{}
	l, err := room.ListRooms()

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	for _, e := range l {
		resp = append(resp, RoomResponse{
			RoomID:      e.RoomID.String(),
			RoomName:    e.Name,
			OwnerName:   e.OwnerName,
			PlayerCount: e.PlayerCount,
			MaxPlayers:  e.MaxPlayers,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (rh *RoomHandler) HandleRoomWebSocket(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	userID := c.QueryParam("user_id")

	hub := rh.roomServer.GetOrCreateHub(roomID, userID)
	err := ws.ServeWs(hub, c.Response(), c.Request(), userID)

	if err != nil {
		return c.String(http.StatusBadRequest, "Error HandleWebSocket")
	}

	return nil
}
