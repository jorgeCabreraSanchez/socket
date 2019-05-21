package model

import (
	"gopkg.in/mgo.v2/bson"
)

type GAuthToken struct {
	UserId bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
}
