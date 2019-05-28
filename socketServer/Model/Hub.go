package model

type Hub struct {
	// set client has participant of a room.
	EnterRoom chan *EnterRoom

	// create chatRoom.
	CreateRoom chan string

	// channels with his participants
	Rooms map[string][]*Client

	// get participant of a room
	UpdatedChatRoom chan string

	//delete user from all rooms where he is
	Unregister chan *Client

	// list of clients that are connected to the socket rigth now
	Clients map[string]*Client

	// set client in clients
	RegisterClient chan *Client
}
