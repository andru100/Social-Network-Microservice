package main

import (
	"context"
	"fmt"
	"net"
	"errors"
	"log"
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

func (s *Server) GetPosts(ctx context.Context, in *model.GetPost) (*model.MongoFields, error) {
	fmt.Println("GetUserComments called!")
	
	switch in.RequestType {
		case "suggested":
			return model.GetSuggestedPosts(ctx, in)
		case "following":
			return model.GetFollowingPosts(ctx, in)
		case "replys":
			return model.GetReplys(ctx, in)
		case "likes":
			return model.GetLikes(ctx, in)
		case "search":
			return model.GetSearchPosts(ctx, in)
		case "user":
			return model.GetUserPosts(ctx, in)
		default:
			return nil, errors.New("Invalid request type")
	}
			
}
