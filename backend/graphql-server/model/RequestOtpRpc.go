package model 

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
)

func RequestOtpRpc (in *RequestOtpInput) (*Confirmation, error) {

	var conn *grpc.ClientConn
	
	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4011", grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewSocialGrpcClient(conn)

	result, err := c.RequestOTP(context.Background(), in)
	
	if err != nil {
		return nil, err
	}
	
	
	return result, err
	
}