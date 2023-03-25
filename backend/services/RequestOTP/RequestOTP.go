package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/andru100/Social-Network-Microservices/backend/services/RequestOTP/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/RequestOTP/utils"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
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

func (s *Server) RequestOTP(ctx context.Context, in *model.RequestOtpInput) (*model.Confirmation, error) { // takes id and sets up bucket and mongodb doc
	fmt.Println("request otp called request type is", in.RequestType)
	
	//create otp
	nums := []rune("123456789")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 6)
	for i := range b {
		b[i] = nums[rand.Intn(len(nums))]
	}

	otp := string(b)

	fmt.Println("randon otp is", otp, "this isnt safe, will need some secret key to truly randomize")

	//save otp to db
	otpHash := utils.HashAndSalt([]byte(otp))

	switch in.RequestType {
	case "sms":

		db := utils.Client.Database("datingapp").Collection("security")

		MobileOTP := &model.Mobile{}

		MobileOTP.Hash = otpHash

		MobileOTP.Expiry = time.Now().Add(time.Minute * 5)

		MobileOTP.Attempts = 0

		filter := bson.M{"Username": in.Username}

		//add new comment to DB
		Updatetype := "$set"
		Key2updt := "OTP.Mobile"
		update := bson.D{
			{Updatetype, bson.D{
				{Key2updt, MobileOTP},
			}},
		}

		//put to db
		_, err := db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, errors.New("its updateone on requestotp 1")
		}

		_, err = model.SendSMS(&in.Mobile, &otp)
		if err != nil {
			return nil, err
		}

		return &model.Confirmation{Username: in.Username, RequestType: in.RequestType}, nil

	case "email":

		fmt.Println("email otp requested")

		db := utils.Client.Database("datingapp").Collection("security")

		EmailOTP := &model.Email{}

		EmailOTP.Hash = otpHash

		EmailOTP.Expiry = time.Now().Add(time.Minute * 5)

		EmailOTP.Attempts = 0

		filter := bson.M{"Username": in.Username}

		//add new comment to DB
		Updatetype := "$set"
		Key2updt := "OTP.Email"
		update := bson.D{
			{Updatetype, bson.D{
				{Key2updt, EmailOTP},
			}},
		}

		//put to db
		_, err := db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, errors.New("its updateone on requestotp2")
		}

		_, err = model.SendEmail(&in.Email, &otp)
		if err != nil {
			return nil, err
		}

		return &model.Confirmation{Username: in.Username, RequestType: in.RequestType}, nil

	case "signup":

		//send sms otp

		fmt.Println("signup otp requested")

		db := utils.Client.Database("datingapp").Collection("tempuser")

		MobileOTP := &model.Mobile{}

		MobileOTP.Hash = otpHash

		MobileOTP.Expiry = time.Now().Add(time.Minute * 5)

		MobileOTP.Attempts = 0

		filter := bson.M{"Username": in.Username}

		//add new comment to DB
		Updatetype := "$set"
		Key2updt := "OTP.Mobile"
		update := bson.D{
			{Updatetype, bson.D{
				{Key2updt, MobileOTP},
			}},
		}

		//put to db
		_, err := db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println("its updateone on requestotp3")
			return nil, errors.New("its updateone on requestotp3")
		}

		//send email otp
		c := make([]rune, 6)
		for i := range b {
			c[i] = nums[rand.Intn(len(nums))]
		}

		otp2 := string(c)

		otpHash = utils.HashAndSalt([]byte(otp2))

		EmailOTP := &model.Email{}

		EmailOTP.Hash = otpHash

		EmailOTP.Expiry = time.Now().Add(time.Minute * 5)

		EmailOTP.Attempts = 0

		Key2updt = "OTP.Email"
		update = bson.D{
			{Updatetype, bson.D{
				{Key2updt, EmailOTP},
			}},
		}

		//put to db
		_, err = db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println("its updateone on requestotp4")
			return nil, errors.New("its updateone on requestotp4")
		}

		_, err = model.SendSMS(&in.Mobile, &otp)
		if err != nil {
			fmt.Println("its send sms on requestotp4", err)
			return nil, err
		}

		_, err = model.SendEmail(&in.Email, &otp2)
		if err != nil {
			fmt.Println("its send email on requestotp4")
			return nil, err
		}

		return &model.Confirmation{Username: in.Username, RequestType: in.RequestType}, nil

	default:
		return nil, fmt.Errorf("Request type not supported")

	}
}
