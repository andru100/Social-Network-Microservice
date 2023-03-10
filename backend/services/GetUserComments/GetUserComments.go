package main

import (
	"context"
	"fmt"
	"net"
	"errors"
	"log"
	"sort"
	"time"
	"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("GetUserComments running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4009))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	model.RegisterSocialGrpcServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (s *Server) GetUserComments(ctx context.Context, in *model.GetComments) (*model.MongoFields, error) {

	
	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.
	currentDoc := model.MongoFields{}
	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	if err != nil {
		err5 := errors.New("unable to find users data")
		fmt.Println(err5, err, in.Username)
		return nil, err5
	}

	sort.Slice(currentDoc.Posts, func(i, j int) bool { // needs to be done on adding post and remove this
		return currentDoc.Posts[i].TimeStamp > currentDoc.Posts[j].TimeStamp
	})

	return &currentDoc, err
}
