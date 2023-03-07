package main

import (
	"context"
	"fmt"
	"net"
	"errors"
	"log"
	//"net/http"
	"sort"
	"time"
	"google.golang.org/grpc"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network-Microservices/GetUserComments/model"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network/backend/social"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
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

	
	collection := social.Client.Database("datingapp").Collection("userdata") // connect to db and collection.
	currentDoc := model.MongoFields{}
	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	if err != nil {
		err = errors.New("unable to find users data")
		return nil, err
	}

	sort.Slice(currentDoc.Posts, func(i, j int) bool { // needs to be done on adding post and remove this
		return currentDoc.Posts[i].TimeStamp > currentDoc.Posts[j].TimeStamp
	})

	return &currentDoc, err
}
