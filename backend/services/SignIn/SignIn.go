package main

import (
	"time"
	"errors"
	"fmt"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"golang.org/x/crypto/bcrypt"
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



func (s *Server) SignIn(ctx context.Context, in *model.UsrsigninInput) (*model.Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	db := utils.Client.Database("datingapp").Collection("security")

	result := model.Security{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	if err != nil {
		return nil, errors.New("username not found")
	}

	fmt.Println("result", result)
	fmt.Println("password", in.Password)

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(in.Password))
    if err != nil {
        log.Println(err)
        return nil, errors.New("password does not match")
    } else {
		token, err1 := model.MakeJwt(&in.Username, true)
		return &model.Jwtdata{Token: token}, err1
	}

}
