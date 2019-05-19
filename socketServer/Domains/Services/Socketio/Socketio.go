package socketio

import (
	"log"
	"socket/socketServer/Domains/Repository/Mongodb"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	gosocketio "github.com/graarh/golang-socketio"
)

func Connection(server *gosocketio.Server, session *mgo.Session) {
	err := server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Print("New Connection")

		// search on db auctions that hi do a bid and subscribe to this room
	})
	if err != nil {
		log.Fatal(err)
	}

}

type Message struct {
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

func Disconnection(server *gosocketio.Server) {
	err := server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Print("Id: " + c.Id() + " Disconnect")
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Error(server *gosocketio.Server) {
	err := server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Println("Error occurs ConnectionId: ", c.Id())
	})
	if err != nil {
		log.Fatal(err)
	}
}

type matchBid struct {
	AuctionId bson.ObjectId `json:"auctionId,omitempty" bson:"auctionId,omitempty"`
	UserId    bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
}

func SubscribeToAuction(server *gosocketio.Server, session *mgo.Session) {
	err := server.On("subscribeToAnAuction", func(c *gosocketio.Channel, object matchBid) {
		_, err := Mongodb.GetBidOfAnAuction(object.AuctionId, object.UserId, session)
		if err == nil {
			c.Join(object.AuctionId.Hex())
			log.Print("subscribed to " + object.AuctionId.Hex())
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
