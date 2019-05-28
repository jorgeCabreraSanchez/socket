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
			EnterRoom:       make(chan map[string]*model.Client),
			CreateRoom:      make(chan string),
			UpdatedChatRoom: make(chan string),
			Unregister:      make(chan *model.Client),
		},
	}

}

func (hubInterface *HubInterface) Run() {
	h := hubInterface.Hub
	for {
		select {
		case enterRoom := <-h.EnterRoom:
			for k, v := range enterRoom {
				h.Rooms[k] = append(h.Rooms[k], v)
			}
		case client := <-h.Unregister:
			for _, clients := range h.Rooms {
				i, err := Helpers.ArrayIndexOf(clients, client)
				if err != nil {
					log.Print("err finding in clients: " + err.Error())
				}
				if i != -1 {
					clients = append(clients[:i], clients[i+1:]...)
					close(client.Send)
				}
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
