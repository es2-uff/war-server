package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	rh *RoomHandler,
) {

	apiRoutes := e.Group("/api/v1")

	apiRoutes.GET("/healthy", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	roomGroup := apiRoutes.Group("/rooms")
	roomGroup.GET("/all", rh.ListRooms)
	roomGroup.POST("/new", rh.CreateNewRoom)
}
