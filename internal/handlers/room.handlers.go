package handlers

import (
	"log"
	"net/http"

	"es2.uff/war-server/internal/domain/room"
	"es2.uff/war-server/internal/ws"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateRoomRequest struct {
	RoomName  string `json:"room_name"`
	OwnerName string `json:"owner_name"`
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

func (rh *RoomHandler) CreateNewRoom(c echo.Context) error {
	r := new(CreateRoomRequest)

	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON format")
	}

	if r.RoomName == "" || r.OwnerName == "" {
		return c.String(http.StatusBadRequest, "Invalid null entry")
	}

	ownerID := uuid.New()
	nr, err := room.NewRoom(r.RoomName, ownerID, r.OwnerName)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	response := RoomResponse{
		RoomID:      nr.RoomID.String(),
		RoomName:    nr.Name,
		OwnerID:     nr.OwnerID.String(),
		OwnerName:   nr.OwnerName,
		PlayerCount: nr.PlayerCount,
		MaxPlayers:  nr.MaxPlayers,
	}

	log.Printf("New room %s created successfully by %s", nr.RoomID, nr.OwnerName)
	return c.JSON(http.StatusOK, response)
}

func (rh *RoomHandler) ListRooms(c echo.Context) error {
	var resp []RoomResponse
	l, err := room.ListRooms()

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	for _, e := range l {
		resp = append(resp, RoomResponse{
			RoomID:      e.RoomID.String(),
			OwnerID:     e.OwnerID.String(),
			RoomName:    e.Name,
			OwnerName:   e.OwnerName,
			PlayerCount: e.PlayerCount,
			MaxPlayers:  e.MaxPlayers,
		})
	}

	if resp == nil {
		resp = []RoomResponse{}
	}

	return c.JSON(http.StatusOK, resp)
}

func (rh *RoomHandler) JoinRoom(roomID uuid.UUID) error {
	panic("Not implemented yet.")
}

func (rh *RoomHandler) HandleWebSocket(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	if roomID == "" {
		return c.String(http.StatusBadRequest, "room_id is required")
	}

	username := c.QueryParam("username")
	if username == "" {
		return c.String(http.StatusBadRequest, "username is required")
	}

	ownerID := c.QueryParam("owner_id")
	userID := uuid.New().String()
	isOwner := false

	if ownerID != "" {
		isOwner = true
		userID = ownerID
	}

	hub := rh.roomServer.GetOrCreateHub(roomID)
	ws.ServeWs(hub, c.Response(), c.Request(), userID, username, isOwner)

	return nil
}
