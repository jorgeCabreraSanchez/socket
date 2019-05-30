package model

type StopListenRoom struct {
	AuctionId string `json:"auctionId,omitempty" bson:"auctionId,omitempty"`
	UserId    string `json:"userId,omitempty" bson:"userId,omitempty"`
}
