package model

type Hub struct {
	// set client has participant of a room.
	EnterRoom chan map[string]*Client

	// create chatRoom.
	CreateRoom chan string

	// channels with his participants
	Rooms map[string][]*Client

	// get participant of a room
	UpdatedChatRoom chan string

	//delete user from all rooms where he is
	Unregister chan *Client
}
