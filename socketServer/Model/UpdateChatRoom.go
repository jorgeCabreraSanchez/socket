package model

import (
	"gopkg.in/mgo.v2/bson"
)

type UpdateChatRoom struct {
	AuctionId bson.ObjectId `json:"auctionId,omitempty" bson:"auctionId,omitempty"`
	Avg       float64       `json:"avg,omitempty" bson:"avg,omitempty"`
}
