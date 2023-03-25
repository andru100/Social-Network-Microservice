package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	//"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func SecurityCheck (in *SecurityCheckInput) (int, error) {

	securityScore := 0

	db := utils.Client.Database("datingapp").Collection("tempuser")

	result := Security{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	if err != nil {
		return securityScore, err
	}

	
	if in.OTP_Mobile != "" {
		fmt.Println("security checking Mobile OTP is", in.OTP_Mobile)
		if result.OTP.Mobile.Expiry.Unix() < time.Now().Unix() {
			return securityScore, errors.New("Mobile OTP expired")
		}

		if result.OTP.Mobile.Attempts > 5 {
			return securityScore, errors.New("too many failed attempts, please request another Mobile OTP")
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.OTP.Mobile.Hash), []byte(in.OTP_Mobile))
		
		if err != nil {
			log.Println(err)

			// log attempt
			result.OTP.Mobile.Attempts += 1

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

			return securityScore, errors.New("mobile OTP does not match")

		} else {

			securityScore += 1

		}
	}

	if in.OTP_Email != "" {
		fmt.Println("security checking Email OTP is", in.OTP_Email)
		if result.OTP.Email.Expiry.Unix() < time.Now().Unix() {
			return securityScore, errors.New("Email OTP expired")
		}

		if result.OTP.Email.Attempts > 5 {
			return securityScore, errors.New("too many failed attempts, please request another Email OTP")
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.OTP.Email.Hash), []byte(in.OTP_Email))
		
		if err != nil {
			log.Println(err)

			// log attempt
			result.OTP.Email.Attempts += 1

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

			return securityScore, errors.New("Email OTP does not match")

		} else {

			securityScore += 1

		}
	}

	if in.Password != "" {
		fmt.Println("security checking password is", in.Password)

		if result.Password.Attempts > 5 {
			return securityScore, errors.New("too many failed password attempts, please reset your password")
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.Password.Hash), []byte(in.Password))
		if err != nil {

			result.Password.Attempts += 1

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
			return securityScore, errors.New("password does not match")
		} 
		
		securityScore += 1
	}

	if in.Token != "" {

		var jwtKey = []byte("AllYourBase")

		claims := &ClaimsChk{}

		tkn, err := jwt.ParseWithClaims(*&in.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			return securityScore, errors.New("JWT token error")
		}
		if !tkn.Valid {
			return securityScore, errors.New("JWT token invalid")
		} 

		securityScore += 1
	}

	return securityScore, err
			
}



	
	

	