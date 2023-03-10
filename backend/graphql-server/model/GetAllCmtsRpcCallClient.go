package model 

import (
	"context"
	"log"
	"google.golang.org/grpc"
)

func Rpc2GetAllCmts (in *GetComments) (*MongoFields, error) {

	var conn *grpc.ClientConn
	
	conn, err := grpc.Dial(":4008", grpc.WithInsecure())
	
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