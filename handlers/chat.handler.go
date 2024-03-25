package handlers

import (
	"chat/models"
	"chat/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func ChatHandler(c echo.Context, as *services.AppService) error {

	ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to upgrade connection")
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("\n %s \n", message)

		var event models.Event
		err = event.Unmarshal(message)
		if err != nil {
			log.Println(err)
			continue
		}

		// TODO: handle error
		switch event.Type {
		case models.JoinEvent:
			_ = as.RoomService.JoinRoom(event.RoomId, event.Username, ws)
		case models.LeaveEvent:
			_ = as.RoomService.LeaveRoom(event.RoomId, event.Username)
		case models.MessageEvent:
			messageEventData := event.Data.(*models.MessageEventData)
			_ = as.RoomService.Broadcast(event.RoomId, event.Username, messageEventData.Text)
		case models.DestoryEvent:
			_ = as.RoomService.RemoveRoom(event.RoomId, event.Username)
		default:
			log.Printf("unknown event type")
		}

	}

}
