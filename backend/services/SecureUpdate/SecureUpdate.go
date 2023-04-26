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
	fmt.Println("secure update called")

	
	if in.UpdateType == "bio" { // bio does not need security check
		//return model.UpdateBio(ctx, in)
		return nil, errors.New("bio update not implemented yet")
	}

	username := in.Username // can be updated while keeping orig for db look up

	dbtype := "security"

	if in.UpdateType == "temp" {//need to give this own field in model next rebuild
		dbtype = "tempuser"
	}

	db := utils.Client.Database("datingapp").Collection(dbtype) // connect to db and collection.

	filter := bson.D{{"Username", in.Username}}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	userdata := model.Security{}

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&userdata)

	if err != nil {
		err = errors.New(fmt.Sprintf("unable to get user data.: %v", err))
		return nil, err
	}

	securityscore, err := model.SecurityCheck(in)

	fmt.Println("security score is", securityscore, "security level is", userdata.SecurityLevel)

	switch in.RequestType {
	case "update":
			fmt.Println("request type is update")

			if securityscore >= userdata.SecurityLevel && err == nil {

				if in.UpdateType == "Password" {
					passwordHash := utils.HashAndSalt([]byte(in.UpdateData))

					passwordHolder := model.Password{Hash: passwordHash, Attempts: 0}

					update := bson.D{
						{"$set", bson.D{
							{"Password", passwordHolder},
						}},
					}

					_, err = db.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						return nil, err
					}

					token, err1 := model.MakeJwt(&in.Username, true)
					return &model.Jwtdata{Token: token}, err1
					
					
				}

				if in.UpdateType == "Mobile" {
					verifyMobile := model.Security{}

					ctxMongo, _ = context.WithTimeout(context.Background(), 15*time.Second)

					err = db.FindOne(ctxMongo, bson.M{"Mobile": in.UpdateData}).Decode(&verifyMobile)

					if err == nil {
						err = errors.New("mobile in use")
						return nil, err
					}

				}

				if in.UpdateType == "Email" {
					verifyEmail := model.Security{}

					err = db.FindOne(ctxMongo, bson.M{"Email": in.UpdateData}).Decode(&verifyEmail)

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

	case "stage2":

		fmt.Println("request tpe is stage 2")

		if err != nil  || securityscore < 1{
			return nil, errors.New(fmt.Sprintf("security check failed, score: %v error %v", securityscore, err))
		}

		

		if securityscore > 0 && securityscore < userdata.SecurityLevel {

			//can edit this to roll through auth pattern and allow custom orders
			_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, Mobile: in.Mobile, RequestType: userdata.AuthType, UserType: "user"})

			if err != nil {
				return nil, errors.New(fmt.Sprintf("error requesting otp on signup stage 2: %v", err))
			}
			
			mobileclue := ""

			if userdata.AuthType == "sms" {
				mobileclue = userdata.Mobile[len(userdata.Mobile)-3:] 
			}
			emailclue := userdata.Email[0:3]
		
			return &model.Jwtdata{Token: "proceed", AuthType: userdata.AuthType, MobClue: mobileclue, EmailClue: emailclue}, nil
		}

		if securityscore >= userdata.SecurityLevel  {
		
			return &model.Jwtdata{Token: "update" /* AuthType: userdata.AuthType, MobClue: mobileclue, EmailClue: emailclue */}, nil

		}

	default:

		fmt.Println("request type is default")

		authoverride := userdata.AuthType
		
		if string(in.RequestType[0]) == "!" { // us ! to ovveride and provide type
			authoverride = string(in.RequestType[1:])
		}
		
		_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Mobile: userdata.Mobile, UserType: in.UpdateType, RequestType: authoverride, Email: userdata.Email})
		
				
		if err != nil {
			return nil, err
		}

		mobileclue := ""

		if userdata.AuthType == "sms" {
			mobileclue = userdata.Mobile[len(userdata.Mobile)-3:] 
		}
		
		emailclue := userdata.Email[0:3]
		a:= &model.Jwtdata{Token: "proceed", AuthType: userdata.AuthType, MobClue: mobileclue, EmailClue: emailclue}

		return a, nil
	}

	return nil, errors.New("invalid request type")

}
