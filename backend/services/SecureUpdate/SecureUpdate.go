package main

import (
	"fmt"
	"errors"
	"context"
	"net"
	"log"

	
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SecureUpdate/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/SecureUpdate/model"
)

func main() {

	fmt.Println("SecureUpdate running!")

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

	db := utils.Client.Database("datingapp").Collection("security")

	securityScore , err := model.SecurityCheck(in)

	if securityScore >= 2 && err == nil {
		filter := bson.M{"Username": in.Username} 
 
		Updatetype := "$set"
		Key2updt := in.UpdateType
		update := bson.D{
			{Updatetype, bson.D{
				{Key2updt, in.UpdateData},
			}},
		}

		//put to db
		_, err = db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}

		token, err1 := model.MakeJwt(&in.Username, true)
		return &model.Jwtdata{Token: token}, err1

	} else {

		return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
	}
}
