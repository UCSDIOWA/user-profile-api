package main

import (
	"context"
	"log"
	"net"

	pb "github.com/UCSDIOWA/user-profile-api/protos"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"google.golang.org/grpc"
)

type server struct{}

type mongo struct {
	Operation *mgo.Collection
}

// DB is a pointer to mongo struct
var DB *mongo

func main() {
	// Host mongo server
	m, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatalf("Could not connect to the MongoDB server: %v", err)
	}
	defer m.Close()

	DB = &mongo{m.DB("tea").C("users")} //change collection to Users

	// Host grpc server
	listen, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("Could not listen on port: %v", err)
	}

	log.Println("Hosting server on", listen.Addr().String())

	s := grpc.NewServer()
	pb.RegisterUserProfileAPIServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) CreatePassword(ctx context.Context, crPsswdReq *pb.CreatePasswordRequest) (*pb.CreatePasswordResponse, error) {
	err := DB.Operation.Update(
		bson.M{"email": (*crPsswdReq).Email},
		bson.M{"$set": bson.M{"password": (*crPsswdReq).Password}})
	if err != nil {
		return &pb.CreatePasswordResponse{Success: false}, err
	}
	return &pb.CreatePasswordResponse{Success: true}, nil
}
