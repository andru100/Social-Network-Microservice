package model 

import (
	"context"
	"log"
	"os"
	"google.golang.org/grpc"
)

func Rpc2GetUserCmts (in *GetComments) (*MongoFields, error) {

	var conn *grpc.ClientConn
	
	conn, err := grpc.Dial(os.Getenv("HOSTIP")+"4009", grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewSocialGrpcClient(conn)

	result, err := c.GetUserComments(context.Background(), in)
	
	if err != nil {
		return nil, err
	}
	
	
	return result, err
	
}