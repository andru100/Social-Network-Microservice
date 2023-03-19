package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	//"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/model"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("SignIn running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4001))
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



func (s *Server) SignIn(ctx context.Context, in *model.SecurityCheckInput) (*model.Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	// check username and password are correct and return security score. 

	switch in.RequestType {
		case "stage1":
			securityScore , err := model.SecurityCheck(in)
			// error will be throw if username or password is incorrect
			if err != nil {
				return nil, err
			}
			if securityScore >= 1 {
				result, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Mobile: in.Mobile, RequestType: "sms"})

				if err != nil {
					return nil, err
				}
				return &model.Jwtdata{Token: "proceed"}, nil
			}
		case "stage2":

			securityScore , err := model.SecurityCheck(in)

			if securityScore >= 2 && err == nil {

				//generate jwt
				token, err1 := model.MakeJwt(&in.Username, true)
				return &model.Jwtdata{Token: token}, err1

			} else {

				return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
			}
		default:
			return nil, errors.New("invalid stage")
	}
	return nil, errors.New("invalid stage")

}
