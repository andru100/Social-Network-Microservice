package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/andru100/Social-Network-Microservices/backend/services/ChkAuth/model"
	jwt "github.com/dgrijalva/jwt-go"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("chkauth running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4007))
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

func (s *Server) Chkauth(ctx context.Context, in *model.JwtdataInput) (*model.Authd, error) {
	
	fmt.Println("running chkauth")

	var jwtKey = []byte("AllYourBase")

	claims := &model.ClaimsChk{}

	tkn, err := jwt.ParseWithClaims(*&in.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	var auth model.Authd

	if err != nil {
		fmt.Println("error parsing token err is: ", err)
		return nil, err
	}
	if !tkn.Valid {
		fmt.Println("token is not valid")
		return nil, err
	}

	auth.AuthdUser = claims.Username

	return &auth, err
}
