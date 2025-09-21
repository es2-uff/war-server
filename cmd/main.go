package main

import (
	"fmt"

	"es2.uff/war-server/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	port := "1323"
	e := echo.New()

	handlers.SetupRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
