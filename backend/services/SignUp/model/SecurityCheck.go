package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	//"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
	//jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func SecurityCheck (in *SecurityCheckInput) (int, error) {

	fmt.Println("running security check")

	securityScore := 0
	

	db := utils.Client.Database("datingapp").Collection("tempuser")

	result := Security{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	if err != nil {
		return securityScore, errors.New("user not found")
	}

	// check if account is locked
	if result.SecurityLock.Status == "Permanent" {
		return securityScore, errors.New("your account has been permanently locked, please contact support")
	}
	if result.SecurityLock.Stage > 2 && result.SecurityLock.Expiry.Unix() > time.Now().Unix() {
		return securityScore , errors.New("account locked, after this try your acccount will be barred and you will need to contact support. Last try in: " + result.SecurityLock.Expiry.Sub(time.Now()).String())
	}
	if result.SecurityLock.Stage > 0 && result.SecurityLock.Expiry.Unix() > time.Now().Unix() {
		return securityScore , errors.New("account locked you can try again in " + result.SecurityLock.Expiry.Sub(time.Now()).String())
	}

	if in.Password != "" {
		fmt.Println("performing seurity check on password", in.Password)

		err = bcrypt.CompareHashAndPassword([]byte(result.Password.Hash), []byte(in.Password))
		if err != nil {

			result.Password.Attempts += 1

			if result.Password.Attempts > 4 {
				// lock account
				stage, err := result.LockAccount()
				
				if err != nil { //if error in locking acount will log this
					return securityScore, err
				}
				return securityScore, errors.New(fmt.Sprintf("Password incorrect, you have exceeded 5 attempts, for security your account has been locked, you can try again in %v minutes, Lock Stage: %v", result.SecurityLock.Expiry.Sub(time.Now()).Minutes(), stage))
			}

			filter := bson.M{"Username": result.Username} 

			Updatetype := "$set"
			Key2updt := "Password"
			update := bson.D{
				{Updatetype, bson.D{
					{Key2updt, result.Password},
				}},
			}

			_, err = db.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return securityScore, err
			}
			fmt.Println("password incorrect, security sore: ", securityScore)
			return securityScore, errors.New("password does not match")
		} 
		
		
		securityScore += 1
		fmt.Println("password correct: security score +1 now is: ", securityScore)
	}

	
	if in.OTP_Mobile != "" {
		fmt.Println("security checking Mobile OTP is", in.OTP_Mobile)
		if result.OTP.Mobile.Expiry.Unix() < time.Now().Unix() {
			return securityScore, errors.New("Mobile OTP expired")
		}


		err = bcrypt.CompareHashAndPassword([]byte(result.OTP.Mobile.Hash), []byte(in.OTP_Mobile))
		
		if err != nil {
			log.Println(err)

			// log attempt
			result.OTP.Mobile.Attempts += 1

			if result.OTP.Mobile.Attempts > 4 {
				// lock account
				stage, err := result.LockAccount()
				if err != nil { //if error in locking acount will log this
					return securityScore, err
				}
				return securityScore, errors.New(fmt.Sprintf("Mobile OTP incorrect, you have exceeded 5 attempts, for security your account has been locked, you can try again in %v minutes, LockStage: %v", result.SecurityLock.Expiry.Sub(time.Now()).Minutes(), stage))
			}

			//put to db

			filter := bson.M{"Username": in.Username} 

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
				return securityScore, err
			}
			fmt.Println("mobileotp incorrect, security score: ", securityScore)
			return securityScore, errors.New("mobile OTP does not match")

		} else {
			securityScore += 1
			fmt.Println("mobileotp correct, security score +1 is: ", securityScore)
			

		}
	}

	if in.OTP_Email != "" {
		fmt.Println("security checking Email OTP is", in.OTP_Email)
		
		if result.OTP.Email.Expiry.Unix() < time.Now().Unix() {
			return securityScore, errors.New("Email OTP expired")
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.OTP.Email.Hash), []byte(in.OTP_Email))
		
		if err != nil {
			log.Println(err)

			// log attempt
			result.OTP.Email.Attempts += 1

			if result.OTP.Email.Attempts > 4 {
				// lock account
				stage, err := result.LockAccount()
				if err != nil { //if error in locking acount will log this
					return securityScore, err
				}
				return securityScore, errors.New(fmt.Sprintf("Email OTP incorrect, you have exceeded 5 attempts, for security your account has been locked, you can try again in %v minutes, LockStage:%v", result.SecurityLock.Expiry.Sub(time.Now()).Minutes(), stage))
			}

			//put to db

			filter := bson.M{"Username": in.Username} 

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
				return securityScore, err
			}
			fmt.Println("emailotp incorrect, security score: ", securityScore)
			return securityScore, errors.New("Email OTP does not match")

		} else {

			securityScore += 1
			fmt.Println("emailotp correct, security score +1 is: ", securityScore)
		}
	}


	fmt.Println("final security score is: ", securityScore)
	return securityScore, err
			
}



	
	

	