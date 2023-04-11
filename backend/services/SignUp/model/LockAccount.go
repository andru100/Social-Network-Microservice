package model

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
)


func (result *Security) LockAccount() (int, error) { // expires otp after 5 minutes


	switch result.SecurityLock.Stage {
	case 0:
		result.SecurityLock.Stage = 1
		result.SecurityLock.Status = "Locked"
		result.SecurityLock.Expiry = time.Now().Add(time.Minute * 1)
	case 1:
		result.SecurityLock.Stage = 2
		result.SecurityLock.Expiry = time.Now().Add(time.Minute * 2)
	case 2:
		result.SecurityLock.Stage = 3
		result.SecurityLock.Expiry = time.Now().Add(time.Minute * 2)
	case 3:
		result.SecurityLock.Status = "Permanent"
	}

	filter := bson.M{"Username": result.Username}

	// update lock account

	Updatetype := "$set"
	Key2updt := "SecurityLock"
	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, result.SecurityLock},
		}},
	}
	
	db := utils.Client.Database("datingapp").Collection("security")

	_, err := db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  result.SecurityLock.Stage , errors.New(fmt.Sprintf("Error loccking account %v", err))
	}

	//reset attempts

	Key2updt = "Password.Attempts"
	update = bson.D{
		{"$set", bson.D{
			{Key2updt, 0},
		}},
	}

	_, err = db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  result.SecurityLock.Stage , errors.New(fmt.Sprintf("Error loccking account %v", err))
	}

	Key2updt = "OTP.Mobile.Attempts"
	update = bson.D{
		{"$set", bson.D{
			{Key2updt, 0},
		}},
	}

	_, err = db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  result.SecurityLock.Stage , errors.New(fmt.Sprintf("Error loccking account %v", err))
	}

	Key2updt = "OTP.Email.Attempts"
	update = bson.D{
		{"$set", bson.D{
			{Key2updt, 0},
		}},
	}

	_, err = db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return  result.SecurityLock.Stage , errors.New(fmt.Sprintf("Error loccking account %v", err))
	}

	return result.SecurityLock.Stage, nil

}