package utils

import (
	"fmt"
	"errors"
	"time"
	"context"
	"log"

	
	//"google.golang.org/grpc"
	"golang.org/x/net/context"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignIn/model"
)

func SecurityCheck (in SecurityCheckInput) (int, error) {

	securityScore := 0

	db := utils.Client.Database("datingapp").Collection("security")

	result := model.Security{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	if err != nil {
		return securityScore, err
	}

	switch {
		case in.OTP_Mobile != "":
			if result.OTP.Mobile.Expiry < time.Now().Unix() {
				return securityScore, errors.New("OTP expired")
			}

			if result.OTP.Mobile.Attempts > 5 {
				return securityScore, errors.New("too many failed attempts, please request another OTP")
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
				_, err = collection.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					return securityScore, err
				}

				return securityScore, errors.New("OTP does not match")

			} else {

				securityScore += 1

			}

		case in.OTP_Email != "":
			if result.OTP.Email.Expiry < time.Now().Unix() {
				return securityScore, errors.New("OTP expired")
			}

			if result.OTP.Email.Attempts > 5 {
				return securityScore, errors.New("too many failed attempts, please request another OTP")
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
				_, err = collection.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					return securityScore, err
				}

				return securityScore, errors.New("OTP does not match")

			} else {

				securityScore += 1

			}
		case in.Password != "": 
			err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(in.Password))
			if err != nil {
				log.Println(err)
				return securityScore, errors.New("password does not match")
			} 
			securityScore += 1

		case in.Token != "":
			fmt.Println("token iss", in.Token)

			var jwtKey = []byte("AllYourBase")

			claims := &model.ClaimsChk{}

			tkn, err := jwt.ParseWithClaims(*&in.Token, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				return securityScore, err
			}
			if !tkn.Valid {
				return securityScore, err
			} 

			securityScore += 1

			
	}



	
	

	
}