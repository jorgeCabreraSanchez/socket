package Hub

import (
	"encoding/json"
	"log"
	"socket/socketServer/Domains/Repository/Mongodb"
	"socket/socketServer/Helpers"
	model "socket/socketServer/Model"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
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
			Clients:         make(map[string]*model.Client),
			EnterRoom:       make(chan *model.EnterRoom),
			CreateRoom:      make(chan string),
			UpdatedChatRoom: make(chan *model.UpdateChatRoom),
			Unregister:      make(chan *model.Client),
			UnregisterRoom:  make(chan string),
			StopListenRoom:  make(chan *model.StopListenRoom),
			RegisterClient:  make(chan *model.Client),
		},
	}

}

func (hubInterface *HubInterface) CreateExistenRooms(db *mgo.Session) {
	auctions, err := Mongodb.GetActualAuctions(db)
	if err != nil {
		log.Print("find first auctions: ", err)
		return
	}

	for _, auction := range auctions {
		hubInterface.Hub.CreateRoom <- auction["_id"].(bson.ObjectId).Hex()
	}
}

func (hubInterface *HubInterface) Run() {
	h := hubInterface.Hub
	for {
		select {
		case enterRoom := <-h.EnterRoom:
			log.Print("enter in room " + enterRoom.AuctionId)
			if client, ok := h.Clients[enterRoom.UserId]; ok {
				h.Rooms[enterRoom.AuctionId] = append(h.Rooms[enterRoom.AuctionId], client)
			}

		case stopListenRoom := <-h.StopListenRoom:
			log.Printf("Client %v stop listen room %v", stopListenRoom.UserId, stopListenRoom.AuctionId)
			i, err := Helpers.ArrayIndexOf(h.Rooms[stopListenRoom.AuctionId], h.Clients[stopListenRoom.UserId])
			if err != nil {
				log.Print("err finding in clients: " + err.Error())
			}
			if i != -1 {
				h.Rooms[stopListenRoom.AuctionId] = append(h.Rooms[stopListenRoom.AuctionId][:i], h.Rooms[stopListenRoom.AuctionId][i+1:]...)
			}

		case client := <-h.Unregister:
			log.Print("unregister")
			// client disconnect from the socket
			// delete him from all rooms
			for k, clients := range h.Rooms {
				i, err := Helpers.ArrayIndexOf(clients, client)
				if err != nil {
					log.Print("err finding in clients: " + err.Error())
				}
				if i != -1 {
					h.Rooms[k] = append(clients[:i], clients[i+1:]...)
				}
			}

			// delete from clients
			if _, ok := h.Clients[client.UserId.Hex()]; ok {
				delete(h.Clients, client.UserId.Hex())
				close(client.Send)
			}

		case roomId := <-h.UnregisterRoom:
			log.Printf("Room %v deleted", roomId)
			if _, ok := h.Rooms[roomId]; ok {
				delete(h.Rooms, roomId)
			}

		case roomId := <-h.CreateRoom:
			log.Print("room " + roomId + " created")
			if _, ok := h.Rooms[roomId]; !ok {
				h.Rooms[roomId] = []*model.Client{}
			}

		case updateChatRoom := <-h.UpdatedChatRoom:
			log.Print("Update room ", updateChatRoom.AuctionId)
			if clients, ok := h.Rooms[updateChatRoom.AuctionId.Hex()]; ok {
				for _, client := range clients {
					avg, _ := json.Marshal(updateChatRoom.Avg)
					client.Send <- avg
				}
			}

		case client := <-h.RegisterClient:
			log.Print("register")
			h.Clients[client.UserId.Hex()] = client
		}

		log.Print("clients: ", h.Clients)
		log.Printf("rooms: %v", h.Rooms)
		log.Print("---------------------")
	}
}
