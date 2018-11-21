package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/UCSDIOWA/user-profile-api/protos"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type server struct{}

type mongo struct {
	Operation *mgo.Collection
}

type getUserProfileResponseStruct struct {
	ProfileImage       string   `json:"profile_image" bson:"profile_image"`
	ProfileDescription string   `json:"profile_description" bson:"profile_description"`
	Endorsements       []string `json:"endorsements" bson:"endorsements"`
	CurrentProjects    []string `json:"current_projects" bson:"current_projects"`
	PreviousProjects   []string `json:"previous_projects" bson:"previous_projects"`
}

// DB is a pointer to mongo struct
var (
	DB           *mongo
	echoEndpoint = flag.String("echo_endpoint", "localhost:50052", "endpoint of user-profile-api")
)

func main() {
	errors := make(chan error)

	go func() {
		errors <- startGRPC()
	}()

	go func() {
		flag.Parse()
		defer glog.Flush()

		errors <- startHTTP()
	}()

	for err := range errors {
		log.Fatal(err)
		return
	}
}

func startGRPC() error {
	// Host mongo server
	m, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatalf("Could not connect to the MongoDB server: %v", err)
	}
	defer m.Close()
	log.Println("Connected to the MongoDB Server.")

	DB = &mongo{m.DB("tea").C("users")}

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

	return err
}

func startHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterUserProfileAPIHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}
	log.Println("Listening on port 8080")

	herokuPort := os.Getenv("PORT")
	if herokuPort == "" {
		herokuPort = "8080"
	}

	return http.ListenAndServe(":"+herokuPort, mux)
}

func (s *server) GetUserProfile(ctx context.Context, request *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	var response pb.GetUserProfileResponse
	var responseStruct getUserProfileResponseStruct
	err := DB.Operation.Find(bson.M{"email": request.Email}).One(&responseStruct)
	if err != nil {
		return nil, err
	}
	response.ProfileImage = responseStruct.ProfileImage
	response.ProfileDescription = responseStruct.ProfileDescription
	response.Endorsements = responseStruct.Endorsements
	response.CurrentProjects = responseStruct.CurrentProjects
	response.PreviousProjects = responseStruct.PreviousProjects

	return &response, nil
}

func (s *server) UpdateUserProfile(ctx context.Context, request *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	find := bson.M{"email": request.Email}
	update := bson.M{"$set": bson.M{"profile_image": request.ProfileImage, "profile_description": request.ProfileDescription, "current_projects": request.CurrentProjects, "previous_projects": request.PreviousProjects}}

	err := DB.Operation.Update(find, update)
	if err != nil {
		return &pb.UpdateUserProfileResponse{Success: false}, nil
	}

	return &pb.UpdateUserProfileResponse{Success: true}, nil
}
