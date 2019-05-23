package model

import (
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	userId bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`

	// The websocket connection.
	conn websocket.Conn `json:"conn,omitempty" bson:"conn,omitempty"`

	// Buffered channel of outbound messages.
	Send chan []byte `json:"send,omitempty" bson:"send,omitempty"`
}
