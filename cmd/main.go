package main

import (
	"fmt"
	"os"

	"es2.uff/war-server/internal/handlers"
	"es2.uff/war-server/internal/ws"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TESTE CD
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize WebSocket room server
	roomServer := ws.NewRoomServer()

	// Initialize handlers
	roomHandler := handlers.NewRoomHandler(roomServer)

	handlers.SetupRoutes(e, roomHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
