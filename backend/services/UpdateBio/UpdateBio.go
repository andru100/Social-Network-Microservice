package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"errors"
	"time"
	"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/UpdateBio/model"
	"github.com/andru100/Social-Network/backend/social"
	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("UpdateBio running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4006))
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


func (s *Server) UpdateBio (ctx context.Context, in *model.UpdateBioInput) (*model.MongoFields, error) { // updates user bio section
	
	collection := social.Client.Database("datingapp").Collection("userdata")

	filter := bson.M{"Username": in.Username}

	Updatetype := "$set"
	Key2updt := "Bio"

	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, in.Bio},
		}},
	}

	//put to db
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err = errors.New("error when updating to DB")
		return nil, err
	}

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err = collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)

	return &currentDoc, err
}
