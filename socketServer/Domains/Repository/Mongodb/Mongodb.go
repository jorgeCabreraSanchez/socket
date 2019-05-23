package Mongodb

import (
	"errors"
	"gnommoApiRest/Config"
	model "socket/socketServer/Model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongoConfig Config.Mongo = Config.GetAll().Mongo

func MongoStart() *mgo.Database {
	session, err :=
		mgo.Dial(mongoConfig.Address + `:` + mongoConfig.Port)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(mongoConfig.Name)
	return db
}

func SetCollection(session *mgo.Session, collection string) *mgo.Collection {
	return session.DB(mongoConfig.Name).C(collection)
}

func GetAuctionsThatIdoABidWithHisAvg(userId bson.ObjectId, db *mgo.Session) ([]bson.M, error) {
	dbsession := db.Copy()
	defer dbsession.Close()
	collection := SetCollection(dbsession, "bids")

	var pipeline []bson.M
	match := bson.M{"userId": userId, "instance.status": "ACTIVE", "won": bson.M{"$exists": false}}
	pipeline = append(pipeline, match)

	lookup := bson.M{"$lookup": bson.M{"from": "bids", "localField": "auctionId", "foreignField": "auctionId", "as": "auctions"}}
	pipeline = append(pipeline, lookup)

	var modelToReturn []bson.M
	pipe := collection.Pipe(pipeline)
	errFind := pipe.All(&modelToReturn)
	if errFind != nil {
		return nil, errFind
	}

	return modelToReturn, nil

}

func GetBidOfAnAuction(auctionId bson.ObjectId, userId bson.ObjectId, db *mgo.Session) (bson.M, error) {
	dbsession := db.Copy()
	defer dbsession.Close()
	collection := SetCollection(dbsession, "bids")

	if auctionId != "" {

		var pipeline []bson.M
		match := bson.M{"auctionId": auctionId, "userId": userId, "instance.status": "ACTIVE", "won": bson.M{"$exists": false}}
		pipeline = append(pipeline, match)

		var modelToReturn bson.M
		pipe := collection.Pipe(pipeline)
		errFind := pipe.One(&modelToReturn)
		if errFind != nil {
			return nil, errFind
		}

		return modelToReturn, nil
	} else {
		return nil, nil
	}
}

func ExistsToken(token string, db *mgo.Session) (model.GAuthToken, error) {
	dbsession := db.Copy()
	defer dbsession.Close()
	collection := SetCollection(dbsession, "gAuthToken")

	var modelToReturn model.GAuthToken

	if token != "" {

		var pipeline []bson.M
		match := bson.M{"$match": bson.M{"token": token}}
		pipeline = append(pipeline, match)

		pipe := collection.Pipe(pipeline)
		errFind := pipe.One(&modelToReturn)
		if errFind != nil {
			return modelToReturn, errFind
		}

		return modelToReturn, nil

	} else {
		return modelToReturn, errors.New("Unauthorized")
	}
}
