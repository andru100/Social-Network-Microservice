package model 

import (
	"context"
	//"fmt"
	"log"
	//"encoding/json"

	//"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/protobuf"

	//"github.com/99designs/gqlgen/graphql"
	//"github.com/andru100/Graphql-Social-Network/graph/GrpcModel"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	//"github.com/andru100/Social-Network/backend/social"
)

func Rpc2GetAllCmts (in *GetComments) (*MongoFields, error) {

	var conn *grpc.ClientConn
	
	conn, err := grpc.Dial(":4009", grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewSocialGrpcClient(conn)

	result, err := c.GetAllComments(context.Background(), in)
	
	if err != nil {
		return nil, err
	}
	
	
	return result, err
	
}