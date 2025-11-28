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
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize WebSocket room server
	roomServer := ws.NewRoomServer()

	// Initialize game manager
	gameManager := ws.NewGameManager()

	// Initialize handlers
	roomHandler := handlers.NewRoomHandler(roomServer)
	gameHandler := handlers.NewGameHandler(gameManager)

	handlers.SetupRoutes(e, roomHandler, gameHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
