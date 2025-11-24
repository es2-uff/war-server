package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	rh *RoomHandler,
	gh *GameHandler,
) {

	apiRoutes := e.Group("/api/v1")

	apiRoutes.GET("/healthy", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	// Player routes
	playerGroup := apiRoutes.Group("/players")
	playerGroup.POST("/new", rh.CreatePlayer)

	// Room routes
	roomGroup := apiRoutes.Group("/rooms")
	roomGroup.GET("/all", rh.ListRooms)
	roomGroup.POST("/new", rh.CreateNewRoom)
	roomGroup.GET("/ws", rh.HandleRoomWebSocket)

	// WebSocket endpoint
	gameGroup := apiRoutes.Group("/games")
	gameGroup.GET("/ws", gh.HandleGameWebSocket)
}
