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
	collection := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// search for duplicate username
	//TODO change this to a map rather than search all docs
	userdata := model.Security{}

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&userdata)

	if err != nil {
		return nil, errors.New("username not found")
	}

	securityscore, err := model.SecurityCheck(in)

	fmt.Println("security score is", securityscore, "security level is", userdata.SecurityLevel)

	switch in.RequestType {
	case "update":
			fmt.Println("update called details are: ", in.Username, in.Password, in.OTP_Mobile, in.OTP_Email, in.UpdateType, in.UpdateData)
			db := utils.Client.Database("datingapp").Collection("security")

			filter := bson.M{"Username": in.Username} 

			ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

			securityscore , err := model.SecurityCheck(in)

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
						fmt.Println("mobile in use error, says it found dupliate, found mob is: ", verifyMobile.Mobile, "full doc is: ", verifyMobile)
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

		fmt.Println("secure update stage 2 called: ", in)

		if err != nil  || securityscore < 1{
			return nil, errors.New(fmt.Sprintf("security check failed, score: %v error %v", securityscore, err))
		}

		

		if securityscore > 0 && securityscore < userdata.SecurityLevel {

			//can edit this to roll through auth pattern and allow custom orders
			_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, Mobile: in.Mobile, RequestType: userdata.AuthType, UserType: "user"})

			if err != nil {
				fmt.Println(err)
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

	// case "stage1":

	// 	fmt.Println("signup stage 1 called: ", in) // always email check fir now can be change


	// 	if err != nil  || securityscore < 1{
	// 		return nil, errors.New(fmt.Sprintf("security check failed, score: %v error %v", securityscore, err))
	// 	}

		

	// 	if securityscore >= 1  {

	// 		db := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

	// 		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// 		// search for duplicate username
	// 		//TODO change this to a map rather than search all docs
	// 		userdata := model.Security{}

	// 		err = db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&userdata)

	// 		if err != nil {
	// 			fmt.Println("error finding user to send otp: ", err)
	// 			err = errors.New("cant find user to send otp")
	// 			return nil, err
	// 		}

			
	// 		_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, Mobile: in.Mobile, RequestType: userdata.AuthType, UserType: "user"})

	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return nil, errors.New(fmt.Sprintf("error requesting otp on signup stage 2: %v", err))
	// 		}
		
			
	// 		mobileclue := ""

	// 		if userdata.AuthType == "sms" {
	// 			mobileclue = userdata.Mobile[len(userdata.Mobile)-3:] 
	// 		}
	// 		emailclue := userdata.Email[0:3]
		
	// 		return &model.Jwtdata{Token: "proceed", AuthType: userdata.AuthType, MobClue: mobileclue, EmailClue: emailclue}, nil
	// 	}

	// 
	default:

		fmt.Println("default called, request is: ", in.RequestType, "username is: ", in.Username)
		dbtype := "security"
		if in.UpdateType == "temp" {//need to give this own field in model next rebuild
			dbtype = "tempuser"
		}
		db := utils.Client.Database("datingapp").Collection(dbtype) // connect to db and collection.

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
		// if in.RequestType == "forgot" {
		// 	if securitydata.AuthType == "none" {
		// 		securitydata.AuthType = "email"
		// 	} else if securitydata.AuthType == "sms" {
		// 		securitydata.AuthType = "both"
		// 	} else if securitydata.AuthType == "email" {
		// 		securitydata.AuthType = "both"
		// 	}
		// } else

		authoverride := securitydata.AuthType
		
		if string(in.RequestType[0]) == "!" { // us ! to ovveride and provide type
			authoverride = string(in.RequestType[1:])
		}
		
		_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Mobile: securitydata.Mobile, UserType: in.UpdateType, RequestType: authoverride, Email: securitydata.Email})
		
				
				
		if err != nil {
			return nil, err
		}

		mobileclue := ""

		if securitydata.AuthType == "sms" {
			mobileclue = securitydata.Mobile[len(securitydata.Mobile)-3:] 
		}
		
		emailclue := securitydata.Email[0:3]
		a:= &model.Jwtdata{Token: "proceed", AuthType: securitydata.AuthType, MobClue: mobileclue, EmailClue: emailclue}
				
		fmt.Println("sending this here:", a)

		return a, nil
	}

	return nil, errors.New("invalid request type")

}
