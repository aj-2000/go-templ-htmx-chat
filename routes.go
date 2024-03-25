package main

import (
	"chat/handlers"
	"chat/services"
	"chat/views/home"

	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo, as *services.AppService) {
	e.GET("chat", func(c echo.Context) error {
		return handlers.ChatHandler(c, as)
	})

	homeComponent := home.HomeIndex("ajay", "aj2000", false, nil, nil)
	e.GET("/", serveComponentHandler(&homeComponent))

	e.GET("/chat", func(c echo.Context) error {
		return handlers.ChatHandler(c, as)
	})
}

func serveComponentHandler(component *templ.Component) echo.HandlerFunc {
	componentHandler := templ.Handler(*component)
	return func(c echo.Context) error {
		return componentHandler.Component.Render(context.TODO(), c.Response().Writer)
	}
}
