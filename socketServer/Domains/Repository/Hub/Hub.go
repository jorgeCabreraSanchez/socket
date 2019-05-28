package Hub

import (
	"encoding/json"
	"log"
	"socket/socketServer/Helpers"
	model "socket/socketServer/Model"
)

// Hub maintains the set of active clients and messages to the
// clients.

type HubInterface struct {
	Hub *model.Hub
}

func NewHub() *HubInterface {
	return &HubInterface{
		Hub: &model.Hub{
			Rooms:           make(map[string][]*model.Client),
			EnterRoom:       make(chan *model.EnterRoom),
			CreateRoom:      make(chan string),
			UpdatedChatRoom: make(chan string),
			Unregister:      make(chan *model.Client),
			Clients:         make(map[string]*model.Client),
			RegisterClient:  make(chan *model.Client),
		},
	}

}

func (hubInterface *HubInterface) Run() {
	h := hubInterface.Hub
	for {
		select {
		case enterRoom := <-h.EnterRoom:
			if client, ok := h.Clients[enterRoom.UserId]; ok {
				h.Rooms[enterRoom.AuctionId] = append(h.Rooms[enterRoom.AuctionId], client)
			}

		case client := <-h.Unregister:
			// client disconnect from the socket
			// delete him from all rooms
			for _, clients := range h.Rooms {
				i, err := Helpers.ArrayIndexOf(clients, client)
				if err != nil {
					log.Print("err finding in clients: " + err.Error())
				}
				if i != -1 {
					clients = append(clients[:i], clients[i+1:]...)
				}
			}

			// delete from clients
			if _, ok := h.Clients[client.UserId.Hex()]; ok {
				delete(h.Clients, client.UserId.Hex())
				close(client.Send)
			}

		case room := <-h.CreateRoom:
			h.Rooms[room] = []*model.Client{}
			log.Print("room " + room + " created")

		case chatRoomId := <-h.UpdatedChatRoom:
			for _, client := range h.Rooms[chatRoomId] {
				tal, _ := json.Marshal("tal")
				// calculate avg and send
				client.Send <- tal
			}

		case client := <-h.RegisterClient:
			// find all bids where he is and enter in his rooms
			h.Clients[client.UserId.Hex()] = client
		}
	}
}
