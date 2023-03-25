package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/utils"
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
				return nil , errors.New(fmt.Sprintf("security check failed: %v", err))
			}
			if securityScore >= 1 {
				collection := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

				ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

				// search for duplicate username
				//TODO change this to a map rather than search all docs
				sendSms := model.Security{}

				err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&sendSms)

				fmt.Println("sendSms.Mobile: ", sendSms.Mobile)

				if err != nil {
					err = errors.New(fmt.Sprintf("unable to locate sms no.: %v", err))
					return nil, err
				}
				_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Mobile: sendSms.Mobile, RequestType: "sms"})

				if err != nil {
					return nil, err
				}
				return &model.Jwtdata{Token: "proceed"}, nil
			}
		case "stage2":

			securityScore , err := model.SecurityCheck(in)

			if securityScore >= 2 && err == nil {

				err = model.UnlockAccount(&in.Username)
				if err != nil {
					return nil, err
				}

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
