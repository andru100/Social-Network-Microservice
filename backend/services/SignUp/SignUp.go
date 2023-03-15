package main

import (
	"context"
	"fmt"
	"log"
	"errors"
	"time"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"golang.org/x/crypto/bcrypt"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("SignUp running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4002))
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



func (s *Server) SignUp(ctx context.Context, newUserData *model.NewUserDataInput) (*model.Jwtdata, error) { // takes id and sets up bucket and mongodb

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// search for duplicate username 
	//TODO change this to a map rather than search all docs
	result := model.MongoFields{}

	err := collection.FindOne(ctxMongo, bson.M{"Username": newUserData.Username}).Decode(&result)

	if err == nil {
		err = errors.New("username in use")
		return nil, err
	}

	createuser := model.MongoFields{Username: newUserData.Username, Email: newUserData.Email, Password: "depreciated", Profpic: "https://adminajh46unique.s3.eu-west-2.amazonaws.com/default-profile-pic.jpg", Photos: []string{}, LastCommentNum: 0, Posts: []*model.PostData{} }

	//username not in use so add new userdata struct
	_, err = collection.InsertOne(context.TODO(), createuser)
	if err != nil {
		return nil, err
	}

	passwordHash := utils.hashAndSalt([]byte(newUserData.Password))
	
	db := utils.Client.Database("datingapp").Collection("security")

	security :=	model.Security{Username: newUserData.Username, Password: passwordHash, OTP: model.OTP{}}

	_ , err = db.InsertOne(context.TODO(), security)
	if err != nil {
		return nil, err
	}


	utils.Createbucket(newUserData.Username) // create bucket to store users files

	//add error return when coial package gets pushed
	token, err2 := model.MakeJwt(&newUserData.Username, true) // make jwt with user id and auth true

	if err2 != nil {
		return nil, err2
	}

	return &model.Jwtdata{Token: token}, err2
}
