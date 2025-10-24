package main

import (
	"fmt"
	"os"

	"es2.uff/war-server/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	port := os.Getenv("PORT")
	e := echo.New()

	roomHandler := handlers.NewRoomHandler()

	handlers.SetupRoutes(e, roomHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
