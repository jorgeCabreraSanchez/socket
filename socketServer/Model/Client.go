package model

import (
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	UserId bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`

	// The websocket connection.
	Conn *websocket.Conn `json:"conn,omitempty" bson:"conn,omitempty"`

	// Buffered channel of outbound messages.
	Send chan []byte `json:"send,omitempty" bson:"send,omitempty"`

	Hub *Hub `json:"hub,omitempty" bson:"hub,omitempty"`

	Unregister chan bool `json:"unregister,omitempty" bson:"unregister,omitempty"`
}
