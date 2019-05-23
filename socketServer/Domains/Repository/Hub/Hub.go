package Hub

import model "socket/socketServer/Model"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// set client has participant of a room.
	enterRoom chan map[string]*model.Client

	// create chatRoom.
	createRoom chan string

	// channels with his participants
	rooms map[string][]*model.Client

	// get participant of a room
	getRooms chan string
}

func newHub() *Hub {
	return &Hub{
		rooms:                  make(map[string][]*model.Client),
		enterRoom:              make(chan map[string]*model.Client),
		createRoom:             make(chan string),
		sendMessageToAChatRoom: make(chan string),
	}
}

func (h *Hub) run() {
	for {
		select {
		case enterRoom := <-h.enterRoom:
			for k, v := range enterRoom {
				h.rooms[k] = append(h.rooms[k], v)
			}
		case room := <-h.createRoom:
			h.rooms[room] = []*model.Client{}
		case chatRoomId := <-h.sendMessageToAChatRoom:
			for room := range h.rooms {

				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
