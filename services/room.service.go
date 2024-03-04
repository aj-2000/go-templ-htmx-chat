package services

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type participant struct {
	username string
	ws       *websocket.Conn
}

type message struct {
	sender string
	text   string
	sentAt time.Time
}

type room struct {
	id           string
	owner        string
	participants map[string]participant
	messages     []message
}

func (r *room) addParticipant(username string, ws *websocket.Conn) error {
	if _, ok := r.participants[username]; ok {
		return errors.New("participant already exists")
	}
	r.participants[username] = participant{username, ws}
	return nil
}

func (r *room) removeParticipant(id string) error {
	if p, ok := r.participants[id]; ok {
		// TODO: handle error
		p.ws.Close()
		delete(r.participants, id)
		return nil
	}
	return errors.New("participant does not exist")
}

func (r *room) broadcast(sender string, text string) {
	m := message{sender, text, time.Now()}
	r.messages = append(r.messages, m)
	for _, p := range r.participants {
		p.ws.WriteJSON(m)
	}
}

func (r *room) destroy() {
	for _, p := range r.participants {
		_ = p.ws.Close()
	}
}

type RoomService struct {
	rooms map[string]*room
}

func NewRoomService() *RoomService {
	return &RoomService{
		rooms: make(map[string]*room),
	}
}

func (rs *RoomService) CreateRoom(roomId string, userId string) error {
	if _, ok := rs.rooms[roomId]; ok {
		return errors.New("room already exists")
	}
	rs.rooms[roomId] = &room{id: roomId, participants: make(map[string]participant), messages: make([]message, 0), owner: userId}
	return nil
}

func (rs *RoomService) RemoveRoom(roomId string, username string) error {
	if _, ok := rs.rooms[roomId]; ok {
		if rs.rooms[roomId].owner != username {
			return errors.New("user is not the owner of the room")
		}
		rs.rooms[roomId].destroy()
		delete(rs.rooms, roomId)
		return nil
	}
	return errors.New("room does not exist")
}

func (rs *RoomService) JoinRoom(roomId string, username string, ws *websocket.Conn) error {
	if r, ok := rs.rooms[roomId]; ok {
		return r.addParticipant(username, ws)
	}
	return errors.New("room does not exist")
}

func (rs *RoomService) LeaveRoom(roomId string, username string) error {
	if r, ok := rs.rooms[roomId]; ok {
		return r.removeParticipant(username)
	}
	return errors.New("room does not exist")
}

func (rs *RoomService) Broadcast(roomId string, username string, text string) error {
	if r, ok := rs.rooms[roomId]; ok {
		r.broadcast(username, text)
		return nil
	}
	return errors.New("room does not exist")
}
