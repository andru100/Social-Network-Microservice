package model

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/utils"
)


func  UnlockAccount(Username *string) (error) { // reset account locks
	
	
	
	sms := "sms"
	email := "email"
	go ResetLockout(Username)
	go ResetPasswordAttempts(Username)
	go ResetMobileOtpAttempts(Username)
	go ResetEmailOtpAttempts(Username)
	go ExpireOTP(Username,  &sms)
	go ExpireOTP(Username,  &email)

	return nil

}

func ResetLockout (Username *string) error {

	resetSecurity := SecurityLock{Status: "Unlocked", Stage: 0,Expiry: time.Now()}

	db := utils.Client.Database("datingapp").Collection("security")
	
	filter := bson.M{"Username": Username}
	update := bson.D{
		{"$set", bson.D{
			{"SecurityLock", resetSecurity},
		}},
	}

	_, err := db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  errors.New(fmt.Sprintf("Error unlocking account %v", err))
	}

	return err
}

func ResetPasswordAttempts (Username *string) error {

	db := utils.Client.Database("datingapp").Collection("security")
	
	filter := bson.M{"Username": Username}
	update := bson.D{
		{"$set", bson.D{
			{"Password.Attempts", 0},
		}},
	}

	_, err := db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  errors.New(fmt.Sprintf("Error resetting password attempts %v", err))
	}

	return err
}

func ResetMobileOtpAttempts (Username *string) error {

	db := utils.Client.Database("datingapp").Collection("security")
	
	filter := bson.M{"Username": Username}
	update := bson.D{
		{"$set", bson.D{
			{ "OTP.Mobile.Attempts", 0},
		}},
	}

	_, err := db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  errors.New(fmt.Sprintf("Error resetting mobile otp attempts %v", err))
	}

	return err
}

func ResetEmailOtpAttempts (Username *string) error {

	db := utils.Client.Database("datingapp").Collection("security")
	
	filter := bson.M{"Username": Username}
	update := bson.D{
		{"$set", bson.D{
			{"OTP.Email.Attempts", 0},
		}},
	}

	_, err := db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  errors.New(fmt.Sprintf("Error resetting email otp attempts %v", err))
	}

	return err
}



func ExpireOTP(Username *string, requestType *string) (error) { // expires otp after 5 minutes
	db := utils.Client.Database("datingapp").Collection("security")
	
	filter := bson.M{"Username": Username}

	switch *requestType {
	case "sms":
		update := bson.D{
			{"$set", bson.D{
				{"OTP.Mobile.Expiry", time.Now()},
			}},
		}
		
		_, err := db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			errors.New(fmt.Sprintf("Error expiring mobile otp for user %v", err))
		}

	case "email":
		update := bson.D{
			{"$set", bson.D{
				{"OTP.Email.Expiry", time.Now()},
			}},
		}
		
		_, err := db.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			errors.New(fmt.Sprintf("Error expiring email otp for user %v", err))
		}

	default:
		errors.New("Invalid request type")

	
	
	}

	return nil
}