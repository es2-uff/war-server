package handlers

import (
	"net/http"

	"es2.uff/war-server/internal/ws"
	"github.com/labstack/echo/v4"
)

type GameHandler struct {
	gameManager *ws.GameManager
}

func NewGameHandler(gameManager *ws.GameManager) *GameHandler {
	return &GameHandler{
		gameManager: gameManager,
	}
}

func (gh *GameHandler) HandleGameWebSocket(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	userID := c.QueryParam("user_id")

	game := gh.gameManager.GetOrCreateGame(roomID)
	err := ws.ServeWs(game, c.Response(), c.Request(), userID)

	if err != nil {
		return c.String(http.StatusBadRequest, "Error HandleWebSocket")
	}

	return nil
}
