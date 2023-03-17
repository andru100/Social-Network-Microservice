package model

import (
	"fmt"
	"math/rand"
	"time"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
)



func (s *Server) RequestOTP (ctx context.Context, in *RequestOtpInput) (*Confirmation, error) {// takes id and sets up bucket and mongodb doc

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

	db := utils.Client.Database("datingapp").Collection("security")

	result := Security{}

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
			result.OTP.Mobile.Hash = otpHash

			result.OTP.Mobile.Expiry = time.Now().Add(time.Minute * 5)

			result.OTP.Mobile.Attempts = 0

			//put to db
			_, err = db.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err := SendSMS(&in.Mobile, &otp)
			if err != nil {
				return nil, err
			}

			return &Confirmation{Username: in.Username, RequestType: in.RequestType}, nil
		
		case "email":

			result.OTP.Email.Hash = otpHash

			result.OTP.Email.Expiry = time.Now().Add(time.Minute * 5)

			result.OTP.Email.Attempts = 0

			//put to db
			_, err = db.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err := SendEmail(&in.Email, &otp)
			if err != nil {
				return nil, err
			}

			return &Confirmation{Username: in.Username, RequestType: in.RequestType}, nil		

		case "signup":

			result.OTP.Mobile.Hash = otpHash

			result.OTP.Mobile.Expiry = time.Now().Add(time.Minute * 5)

			result.OTP.Mobile.Attempts = 0

			//add email otp
			c := make([]rune, 6)
			for i := range b {
				c[i] = nums[rand.Intn(len(nums))]
			}

			otp2 := string(c) 

			//save otp to db
			EmailHash := utils.HashAndSalt([]byte(otp2))

			result.OTP.Email.Hash = EmailHash

			result.OTP.Email.Expiry = time.Now().Add(time.Minute * 5)

			result.OTP.Email.Attempts = 0

			//put to db
			_, err = db.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err = SendSMS(&in.Mobile, &otp)
			if err != nil {
				return nil, err
			}

			_, err = SendEmail(&in.Email, &otp)
			if err != nil {
				return nil, err
			}

			return &Confirmation{Username: in.Username, RequestType: in.RequestType}, nil	
	
		default:
			return nil, fmt.Errorf("Request type not supported")

	
	}
}