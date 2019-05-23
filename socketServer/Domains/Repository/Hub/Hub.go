package Hub

import (
	"encoding/json"
	"log"
	model "socket/socketServer/Model"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// set client has participant of a room.
	EnterRoom chan map[string]*model.Client

	// create chatRoom.
	CreateRoom chan string

	// channels with his participants
	Rooms map[string][]*model.Client

	// get participant of a room
	UpdatedChatRoom chan string
}

func NewHub() *Hub {
	return &Hub{
		Rooms:           make(map[string][]*model.Client),
		EnterRoom:       make(chan map[string]*model.Client),
		CreateRoom:      make(chan string),
		UpdatedChatRoom: make(chan string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case enterRoom := <-h.EnterRoom:
			for k, v := range enterRoom {
				h.Rooms[k] = append(h.Rooms[k], v)
			}
		case room := <-h.CreateRoom:
			h.Rooms[room] = []*model.Client{}
			log.Print("room " + room + " created")
		case chatRoomId := <-h.UpdatedChatRoom:
			for _, client := range h.Rooms[chatRoomId] {
				tal, _ := json.Marshal("tal")
				select {
				case client.Send <- tal:
				default:
					close(client.Send)
					// delete(h.clients, client)
				}
			}
		}
	}
}
