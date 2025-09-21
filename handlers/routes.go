package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
) {
	e.GET("/healthy", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})
}
