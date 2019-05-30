package model

type Hub struct {
	// set client has participant of a room.
	EnterRoom chan *EnterRoom `json:"enterRoom,omitempty" bson:"enterRoom,omitempty"`

	// delete client from array of clients of a room
	StopListenRoom chan *StopListenRoom `json:"stopListenRoom,omitempty" bson:"stopListenRoom,omitempty"`

	// create chatRoom.
	CreateRoom chan string `json:"createRoom,omitempty" bson:"createRoom,omitempty"`

	// channels with his participants
	Rooms map[string][]*Client `json:"rooms,omitempty" bson:"rooms,omitempty"`

	// get participant of a room
	UpdatedChatRoom chan *UpdateChatRoom `json:"updatedChatRoom,omitempty" bson:"updatedChatRoom,omitempty"`

	//delete user from all rooms where he is
	Unregister chan *Client `json:"unregister,omitempty" bson:"unregister,omitempty"`

	// delete room
	UnregisterRoom chan string `json:"unregisterRoom,omitempty" bson:"unregisterRoom,omitempty"`

	// list of clients that are connected to the socket rigth now
	Clients map[string]*Client `json:"clients,omitempty" bson:"clients,omitempty"`

	// set client in clients
	RegisterClient chan *Client `json:"registerClient,omitempty" bson:"registerClient,omitempty"`
}
