package main

import (
	"fmt"
	"math/rand"
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

	fmt.Println("RequestOTP running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4011))
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

func (s *Server) RequestOTP (ctx context.Context, in *model.RequestOtpInput) (*model.Confirmation, error) {// takes id and sets up bucket and mongodb doc

	//create otp
	nums := []rune("123456789")

	rand.Seed(time.Now().UnixNano())

    b := make([]rune, 6)
    for i := range b {
        b[i] = nums[rand.Intn(len(nums))]
    }

	otp := string(b) 

	fmt.Println("randon otp is", otp, "this isnt safe, wiill need some secret key to truly randomize")

	//save otp to db
	passwordHash := utils.HashAndSalt([]byte(otp))

	db := utils.Client.Database("datingapp").Collection("security")

	result := model.Security{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	if err != nil {
		return nil, err
	}

	result.OTP[in.RequestType].Hash = passwordHash

	result.OTP[in.RequestType].Expiry = time.Now().Add(time.Minute * 5)

	result.OTP[in.RequestType].Attempts = 0

	filter := bson.M{"Username": in.Username} 

	//add new comment to DB 
	Updatetype := "$set"
	Key2updt := "OTP"
	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, result.OTP},
		}},
	}

	//put to db
	_, err = db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	//send otp
	if in.requestType == "sms" {
		_, err := utils.SendSMS(in.Mobile, otp)
		if err != nil {
			return nil, err
		}

		return &model.Confirmation{Username: in.UserName, RequestType: in.RequestType}, nil

	} else if in.requestType == "email" {
		_, err := utils.SendEmail(in.Email, otp)
		if err != nil {
			return nil, err
		}

		return &model.Confirmation{Username: in.UserName, RequestType: in.RequestType}, nil		
	}


}