package main

import (
	"fmt"
	"log"
	"net/http"
	"socket/socketServer/Config"
	"socket/socketServer/Domains/Repository/Hub"
	"socket/socketServer/Domains/Repository/Mongodb"
	"socket/socketServer/Domains/Services/Auth"
	socket "socket/socketServer/Domains/Services/Socket"

	"github.com/gorilla/websocket"
)

func main() {
	config := Config.GetAll()
	db := Mongodb.MongoStart()

	fmt.Println("Socket Started")

	var upgrader = websocket.Upgrader{} // use default options

	// socketio.Connection(server, db.Session)
	// socketio.Disconnection(server)
	// socketio.Error(server)
	// socketio.SubscribeToAuction(server, db.Session

	hub := Hub.NewHub()
	go hub.Run()
	http.Handle("/socket.io/", Auth.AuthMiddleware(socket.StartServer(upgrader, hub, db.Session), db.Session))
	log.Fatal(http.ListenAndServe(":"+config.StatusMicro.Port, nil))
}
