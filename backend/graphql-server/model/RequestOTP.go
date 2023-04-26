package model

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) RequestOTP(ctx context.Context, in *RequestOtpInput) (*Confirmation, error) { // takes id and sets up bucket and mongodb doc
	fmt.Println("request otp called")
	//create otp
	nums := []rune("123456789")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 6)
	for i := range b {
		b[i] = nums[rand.Intn(len(nums))]
	}

	otp := string(b)

	fmt.Println("Debug mode: OTP created is: ", otp)

	//save otp to db
	otpHash := utils.HashAndSalt([]byte(otp))

	switch in.RequestType {
	case "sms":
		fmt.Println("sms otp requested")

		db := utils.Client.Database("datingapp").Collection("security")

		MobileOTP := &Mobile{}

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

		_, err = SendSMS(&in.Mobile, &otp)
		if err != nil {
			return nil, err
		}

		return &Confirmation{Username: in.Username, RequestType: in.RequestType}, nil

	case "email":

		fmt.Println("email otp requested")

		db := utils.Client.Database("datingapp").Collection("security")

		EmailOTP := &Email{}

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

		_, err = SendEmail(&in.Email, &otp)
		if err != nil {
			return nil, err
		}

		return &Confirmation{Username: in.Username, RequestType: in.RequestType}, nil

	case "signup":

		//send sms otp

		db := utils.Client.Database("datingapp").Collection("tempuser")

		MobileOTP := &Mobile{}

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
			return nil, errors.New("error updating mobile otp to db")
		}

		//send email otp
		c := make([]rune, 6)
		for i := range b {
			c[i] = nums[rand.Intn(len(nums))]
		}

		otp2 := string(c)

		EmailOTP := &Email{}

		EmailOTP.Hash = otp2

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
			return nil, errors.New("error updating email otp to db")
		}

		_, err = SendSMS(&in.Mobile, &otp)
		if err != nil {
			return nil, err
		}

		_, err = SendEmail(&in.Email, &otp2)
		if err != nil {
			return nil, err
		}

		return &Confirmation{Username: in.Username, RequestType: in.RequestType}, nil

	default:
		return nil, fmt.Errorf("Request type not supported")

	}
}
