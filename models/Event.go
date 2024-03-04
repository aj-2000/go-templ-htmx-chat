package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type EventType string

const (
	JoinEvent    EventType = "join"
	LeaveEvent   EventType = "leave"
	MessageEvent EventType = "message"
	DestoryEvent EventType = "destroy"
)

type Event struct {
	Username string    `json:"username"`
	RoomId   string    `json:"roomId"`
	Type     EventType `json:"type"`
	Data     any       `json:"data"`
}

type JoinEventData struct{}

type LeaveEventData struct{}

type MessageEventData struct {
	Text string `json:"text"`
}

type DestroyEventData struct{}

func (e *Event) Unmarshal(data []byte) error {
	println(string(data))
	var raw map[string]any
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	eventType, ok := raw["type"].(string)
	if !ok {
		return errors.New("invalid event type")
	}
	e.Type = EventType(eventType)

	username, ok := raw["username"].(string)
	if !ok {
		return errors.New("username is req")
	}
	e.Username = username

	roomId, ok := raw["roomId"].(string)
	if !ok {
		return errors.New("roomId is req")
	}
	e.RoomId = roomId

	eventData, ok := raw["data"].(map[string]any)
	if !ok {
		eventData = make(map[string]any)
	}
	switch e.Type {
	case JoinEvent:
		e.Data = &JoinEventData{}
	case LeaveEvent:
		e.Data = &LeaveEventData{}
	case MessageEvent:
		e.Data = &MessageEventData{
			Text: getString(eventData["text"]),
		}
	case DestoryEvent:
		e.Data = &DestroyEventData{}
	default:
		return fmt.Errorf("unknown event type: %v", e.Type)
	}

	if err := validateEventData(e.Data); err != nil {
		return err
	}

	return nil
}

func (e *Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func getString(v any) string {
	if v == nil {
		return ""
	}
	return v.(string)
}

func validateEventData(data any) error {

	switch d := data.(type) {
	case *JoinEventData:
	case *LeaveEventData:
	case *MessageEventData:
		if d.Text == "" {
			return errors.New("invalid message event data")
		}
	case *DestroyEventData:
	default:
		return errors.New("unknown event data type a")
	}
	return nil
}
