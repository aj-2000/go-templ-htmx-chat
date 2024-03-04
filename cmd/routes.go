package main

import (
	"chat/handlers"
	"chat/services"

	"github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo, as *services.AppService) {
	e.GET("/chat", func(c echo.Context) error {
		return handlers.ChatHandler(c, as)
	})
}
