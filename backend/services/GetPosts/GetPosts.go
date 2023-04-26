package main

import (
	"context"
	"fmt"
	"net"
	"errors"
	"log"
	"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/model"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("GetPosts running!")

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
	fmt.Println("get posts service called!")
	
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
