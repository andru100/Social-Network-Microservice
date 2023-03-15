package main

import (
	"fmt"
	"errors"
	"time"
	"context"
	"net"
	"log"

	
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/model"
)

func main() {

	fmt.Println("ConfirmOTP running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4012))
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

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func (s *Server) SecureUpdate (ctx context.Context, in *model.SecurityCheckInput) (*model.Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	
	securityScore , err := utils.SecurityCheck(in)

	if securityScore >= 2 && err == nil {
		filter := bson.M{"Username": in.Username} 
 
		Updatetype := "$set"
		Key2updt := in.UpdateType
		update := bson.D{
			{Updatetype, bson.D{
				{Key2updt, result.OTP},
			}},
		}

		//put to db
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}

		token, err1 := model.MakeJwt(&in.Username, true)
		return &model.Jwtdata{Token: token}, err1

	} else {

		return nil, errors.New("security check failed: %v", err)
	}
}
