package handlers

import (
	"log"
	"net/http"

	"es2.uff/war-server/internal/domain/room"
	"es2.uff/war-server/internal/ws"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateRoomRequest struct { // TODO: Move it to a proper DTO package later
	RoomName string `json:"room_name"`
}

type RoomResponse struct { // TODO: Move it to a proper DTO package later
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
	OwnerID  string `json:"owner_id"`
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

	if r.RoomName == "" {
		return c.String(http.StatusBadRequest, "Invalid null entry")
	}

	nr, err := room.NewRoom(r.RoomName)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	log.Printf("New room %s created succesfully", nr.RoomID)
	return c.JSON(http.StatusOK, nr)
}

func (rh *RoomHandler) ListRooms(c echo.Context) error {
	var resp []RoomResponse
	l, err := room.ListRooms()

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}

	if len(l) == 0 {
		return c.NoContent(http.StatusOK)
	}

	for _, e := range l {
		resp = append(resp, RoomResponse{
			RoomID:   e.RoomID.String(),
			OwnerID:  e.OwnerID.String(),
			RoomName: e.Name,
		})
	}

	return c.JSON(http.StatusOK, l)
}

func (rh *RoomHandler) JoinRoom(roomID uuid.UUID) error {
	panic("Not implemented yet.")
}

func (rh *RoomHandler) HandleWebSocket(c echo.Context) error {
	// Get room ID from query parameter or path parameter
	roomID := c.QueryParam("room_id")
	if roomID == "" {
		roomID = c.Param("room_id")
	}

	if roomID == "" {
		return c.String(http.StatusBadRequest, "room_id is required")
	}

	// Get or create hub for this room
	hub := rh.roomServer.GetOrCreateHub(roomID)

	// Generate a unique user ID for this connection
	userID := uuid.New().String()

	// Upgrade connection to WebSocket and serve
	ws.ServeWs(hub, c.Response(), c.Request(), userID)

	return nil
}
