package main

import (
	"fmt"
	"math/rand"
	"time"
	"context"
	"net"
	"log"

	
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/RequestOTP/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/RequestOTP/model"
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

	filter := bson.M{"Username": in.Username} 

	//add new comment to DB 
	Updatetype := "$set"
	Key2updt := "OTP"
	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, result.OTP},
		}},
	}

	switch in.RequestType {
		case "sms":
			result.OTP.Mobile.Hash = passwordHash

			result.OTP.Mobile.Expiry = time.Now().Add(time.Minute * 5)

			result.OTP.Mobile.Attempts = 0

			//put to db
			_, err = db.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err := model.SendSMS(&in.Mobile, &otp)
			if err != nil {
				return nil, err
			}

			return &model.Confirmation{Username: in.Username, RequestType: in.RequestType}, nil
		
		case "email":
			result.OTP.Email.Hash = passwordHash

			result.OTP.Email.Expiry = time.Now().Add(time.Minute * 5)

			result.OTP.Email.Attempts = 0

			//put to db
			_, err = db.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err := model.SendEmail(&in.Email, &otp)
			if err != nil {
				return nil, err
			}

			return &model.Confirmation{Username: in.Username, RequestType: in.RequestType}, nil		
	
		default:
			return nil, fmt.Errorf("Request type not supported")

	
	}
}