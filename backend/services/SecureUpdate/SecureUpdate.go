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

	


	//todo add updatetype to input on next model rebuild and re right this
	updateData := make(map[string]string)
	updateData["Password"] = "sms+email"
	updateData["Email"] = "sms"
	updateData["Mobile"] = "email"


	 

	switch in.RequestType {
	case "stage2":
	
			db := utils.Client.Database("datingapp").Collection("security")

			securityScore , err := model.SecurityCheck(in)

			if securityScore >= 2 && err == nil {
				filter := bson.M{"Username": in.Username} 
		
				Updatetype := "$set"
				Key2updt := in.RequestType
				update := bson.D{
					{Updatetype, bson.D{
						{Key2updt, in.UpdateData},
					}},
				}

				//put to db
				_, err = db.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					return nil, err
				}

				token, err1 := model.MakeJwt(&in.Username, true)
				return &model.Jwtdata{Token: token}, err1

			} else {

				return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
			}

	default:
		db := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

		sendotp := model.Security{}

		err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&sendotp)

		fmt.Println("sendSms.Mobile: ", sendotp.Mobile)

		if err != nil {
			err = errors.New(fmt.Sprintf("unable to locate sms no.: %v", err))
			return nil, err
		}
		_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Mobile: sendotp.Mobile, RequestType: updateData[in.RequestType]})

		if err != nil {
			return nil, err
		}
		return &model.Jwtdata{Token: "proceed"}, nil
	}

}
