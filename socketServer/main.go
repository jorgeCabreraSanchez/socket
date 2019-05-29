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

	"google.golang.org/grpc"
)

type server struct {
	Hub *model.Hub
}

// SayHello implements helloworld.SocketServer
func (c *server) UploadAuction(ctx context.Context, in *pb.UploadAuctionBody) (*pb.Empty, error) {
	log.Printf("Received: %v", in.AuctionId)
	c.Hub.UpdatedChatRoom <- in.AuctionId
	return &pb.Empty{}, nil
}

func (c *server) ListenRoom(ctx context.Context, in *pb.ListenRoomBody) (*pb.Empty, error) {
	log.Printf("Received: %v", in.AuctionId)
	c.Hub.EnterRoom <- &model.EnterRoom{AuctionId: in.AuctionId, UserId: in.UserId}
	return &pb.Empty{}, nil
}

func main() {
	config := Config.GetAll()
	db := Mongodb.MongoStart()

	// socketio.Connection(server, db.Session)
	// socketio.Disconnection(server)
	// socketio.Error(server)
	// socketio.SubscribeToAuction(server, db.Session

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
	pb.RegisterSocketServer(s, &server{Hub: hub.Hub})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
