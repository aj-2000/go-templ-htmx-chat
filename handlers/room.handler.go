package handlers

import (
	"chat/services"
	"chat/views/chatroom"
	"context"

	"github.com/labstack/echo/v4"
)

func CreateRoomHandler(c echo.Context, as *services.AppService) error {
	return chatroom.ChatRoom().Render(context.TODO(), c.Response().Writer)
}
