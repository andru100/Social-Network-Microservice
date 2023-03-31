package main

import (
	"fmt"
	"errors"
	"context"
	"net"
	"log"
	"time"

	
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SecureUpdate/utils"
	"github.com/andru100/Social-Network-Microservices/backend/services/SecureUpdate/model"
)

func main() {

	fmt.Println("SecureUpdate running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4012))
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

func (s *Server) SecureUpdate (ctx context.Context, in *model.SecurityCheckInput) (*model.Jwtdata, error) {// takes id and sets up bucket and mongodb doc
	fmt.Println("secure update called request type is", in.RequestType)

	username := in.Username // can be updated while keeping orig for db look up
	switch in.RequestType {
	case "update":
			fmt.Println("update called details are: ", in.Username, in.Password, in.OTP_Mobile, in.OTP_Email, in.UpdateType, in.UpdateData)
			db := utils.Client.Database("datingapp").Collection("security")

			filter := bson.M{"Username": in.Username} 

			ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

			

			securityScore , err := model.SecurityCheck(in)

			if securityScore >= 2 && err == nil {

				if in.UpdateType == "Password" {
					passwordHash := utils.HashAndSalt([]byte(in.UpdateData))

					passwordHolder := model.Password{Hash: passwordHash, Attempts: 0}

					update := bson.D{
						{"$set", bson.D{
							{"Password", passwordHolder},
						}},
					}

					//put to db
					_, err = db.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						return nil, err
					}

					token, err1 := model.MakeJwt(&in.Username, true)
					return &model.Jwtdata{Token: token}, err1
					
					
				}

				if in.UpdateType == "Mobile" {
					verifyMobile := model.Security{}

					err = db.FindOne(ctxMongo, bson.M{"Mobile": in.Email}).Decode(&verifyMobile)

					if err == nil {
						err = errors.New("mobile in use")
						return nil, err
					}

				}

				if in.UpdateType == "Email" {
					verifyEmail := model.Security{}

					err = db.FindOne(ctxMongo, bson.M{"Email": in.Email}).Decode(&verifyEmail)

					if err == nil {
						err = errors.New("email in use")
						return nil, err
					}
				}

				

				if in.UpdateType == "Username" {
				
					verifyUsername := model.Security{}

					err := db.FindOne(ctxMongo, bson.M{"Username": in.UpdateData}).Decode(&verifyUsername)

					if err == nil {
						err = errors.New("username in use")
						return nil, err
					}

					//add to userdata as well as security
					userdb := utils.Client.Database("datingapp").Collection("userdata")

					filter := bson.M{"Username": in.Username} 

					update := bson.D{
						{"$set", bson.D{
							{in.UpdateType, in.UpdateData},
						}},
					}
	
					//put to db
					_, err = userdb.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						return nil, err
					}

					
					username = in.UpdateData

				}


				filter := bson.M{"Username": in.Username} 
		
		
				update := bson.D{
					{"$set", bson.D{
						{in.UpdateType, in.UpdateData},
					}},
				}

				//put to db
				_, err = db.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					return nil, err
				}

				token, err1 := model.MakeJwt(&username, true)
				return &model.Jwtdata{Token: token}, err1

			} else {

				return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
			}

		default:

		fmt.Println("default called, reuest is: ", in.RequestType)
		db := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

		securitydata := model.Security{}

		err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&securitydata)

		if err != nil {
			err = errors.New(fmt.Sprintf("unable to get user data.: %v", err))
			return nil, err
		}

		fmt.Println("security data retrieved is: ", securitydata.AuthType, securitydata)
		//overide auth types
		//ifs forgot password, chech auth type and add extra layer eg:none = email, sms = both, email = both
		//in future an add other auths eg auth app, and add to this
		if in.RequestType == "forgot" {
			if securitydata.AuthType == "none" {
				securitydata.AuthType = "email"
			} else if securitydata.AuthType == "sms" {
				securitydata.AuthType = "both"
			} else if securitydata.AuthType == "email" {
				securitydata.AuthType = "both"
			}
		} else if string(in.RequestType[0]) == "!" { // us ! to ovveride and provide type
			securitydata.AuthType = string(in.RequestType[1:])
		}
		
		_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Mobile: securitydata.Mobile, UserType: in.UpdateType, RequestType: securitydata.AuthType, Email: securitydata.Email})
		
				
				
		if err != nil {
			return nil, err
		}

		mobileclue := securitydata.Mobile[len(securitydata.Mobile)-3:] 
		emailclue := securitydata.Email[0:3]
		a:= &model.Jwtdata{Token: "proceed", AuthType: securitydata.AuthType, MobClue: mobileclue, EmailClue: emailclue}
				
		fmt.Println("sending this here:", a)

		return a, nil
	}

	return nil, errors.New("invalid request type")

}
