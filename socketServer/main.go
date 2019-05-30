package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"socket/socketServer/Config"
	"socket/socketServer/Domains/Client"
	"socket/socketServer/Domains/Repository/Hub"
	"socket/socketServer/Domains/Repository/Mongodb"
	"socket/socketServer/Domains/Services/Auth"
	model "socket/socketServer/Model"
	pb "socket/socketServer/proto"

	"gopkg.in/mgo.v2/bson"

	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
)

type server struct {
	Hub *model.Hub
	Db  *mgo.Session
}

// SayHello implements helloworld.SocketServer
func (c *server) UploadAuction(ctx context.Context, in *pb.UploadAuctionBody) (*pb.Empty, error) {
	log.Printf("UploadAuction: %v", in.AuctionId)

	dbsession := c.Db.Copy()
	defer dbsession.Close()

	updateChatRoom, err := Mongodb.GetAvgOfAnAuction(bson.ObjectIdHex(in.AuctionId), dbsession)
	if err != nil {
		log.Print("err uploading: ", err)
		return &pb.Empty{}, err
	}
	c.Hub.UpdatedChatRoom <- &updateChatRoom
	return &pb.Empty{}, nil
}

func (c *server) ListenRoom(ctx context.Context, in *pb.ListenRoomBody) (*pb.Empty, error) {
	log.Printf("ListenRoom: %v", in.AuctionId)
	c.Hub.EnterRoom <- &model.EnterRoom{AuctionId: in.AuctionId, UserId: in.UserId}
	return &pb.Empty{}, nil
}

func (c *server) UnregisterRoom(ctx context.Context, in *pb.UnregisterRoomBody) (*pb.Empty, error) {
	log.Printf("UnregisterRoom: %v", in.AuctionId)
	c.Hub.UnregisterRoom <- in.AuctionId
	return &pb.Empty{}, nil
}

func (c *server) StopListenRoom(ctx context.Context, in *pb.StopListenRoomBody) (*pb.Empty, error) {
	log.Printf("StopListenRoom: %v", in.AuctionId)
	c.Hub.StopListenRoom <- &model.StopListenRoom{AuctionId: in.AuctionId, UserId: in.UserId}
	return &pb.Empty{}, nil
}

func (c *server) CreateRoom(ctx context.Context, in *pb.CreateRoomBody) (*pb.Empty, error) {
	log.Printf("CreateRoom: %v", in.AuctionId)
	c.Hub.CreateRoom <- in.AuctionId
	return &pb.Empty{}, nil
}

func main() {
	config := Config.GetAll()
	db := Mongodb.MongoStart()

	hub := Hub.NewHub()
	go hub.Run()
	hub.CreateExistenRooms(db.Session)
	http.Handle("/socket.io/", Auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Client.ServeWs(hub.Hub, db.Session, w, r)
	}), db.Session))

	go http.ListenAndServe(":"+config.StatusMicro.Port, nil)

	lis, err := net.Listen("tcp", ":"+config.GrpcMicro.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSocketServer(s, &server{Hub: hub.Hub, Db: db.Session})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
