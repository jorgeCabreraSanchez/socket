package main

import (
	"fmt"
	"log"
	"net/http"
	"socket/socketServer/Config"
	"socket/socketServer/Domains/Repository/Mongodb"
	"socket/socketServer/Domains/Services/Auth"
	socketio "socket/socketServer/Domains/Services/Socketio"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func main() {
	config := Config.GetAll()
	db := Mongodb.MongoStart()

	fmt.Println("Socket Started")

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	socketio.Connection(server, db.Session)
	socketio.Disconnection(server)
	socketio.Error(server)
	socketio.SubscribeToAuction(server, db.Session)

	http.Handle("/socket.io/", Auth.AuthMiddleware(server, db.Session))
	log.Fatal(http.ListenAndServe(":"+config.StatusMicro.Port, nil))
}
