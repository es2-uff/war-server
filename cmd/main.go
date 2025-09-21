package main

import (
	"fmt"

	"es2.uff/war-server/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	port := "1323"
	e := echo.New()

	roomHandler := handlers.NewRoomHandler()

	handlers.SetupRoutes(e, roomHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
