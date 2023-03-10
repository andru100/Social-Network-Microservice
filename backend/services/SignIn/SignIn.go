package main

import (
	//"context"
	"time"
	"errors"
	"fmt"
	"log"
	"net"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	//"github.com/andru100/Social-Network-Microservices/Social" once pushed

	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/model"
	//jwt "github.com/dgrijalva/jwt-go"
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



func (s *Server) SignIn(ctx context.Context, in *model.UsrsigninInput) (*model.Jwtdata, error) {// takes id and sets up bucket and mongodb

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection

	result := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	fmt.Println("result", result)
	fmt.Println("password", in.Password)
	
	if err != nil {
		return nil, errors.New("username not found")
	}

	if result.Password == in.Password {
		token, err2 := model.MakeJwt(&in.Username, true)
		return &model.Jwtdata{Token: token}, err2
	} else {
		return nil, errors.New("password does not match")
	}

}
