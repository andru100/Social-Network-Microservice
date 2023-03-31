package model

import (
	"fmt"
	"errors"
	"context"
	"time"

	
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
)


func (s *Server) SecureUpdate (ctx context.Context, in *SecurityCheckInput) (*Jwtdata, error) {// takes id and sets up bucket and mongodb doc
	fmt.Println("secure update called request type is", in.RequestType)

	switch in.RequestType {
	case "update":
	
			db := utils.Client.Database("datingapp").Collection("security")

			securityScore , err := SecurityCheck(in)

			if securityScore >= 2 && err == nil {
				filter := bson.M{"Username": in.Username} 
		
				Updatetype := "$set"
				Key2updt := in.UpdateType
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

				token, err1 := MakeJwt(&in.Username, true)
				return &Jwtdata{Token: token}, err1

			} else {

				return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
			}

		default:
		db := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

		sendotp := Security{}

		err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&sendotp)

		fmt.Println("sendSms.Mobile: ", sendotp.Mobile)

		if err != nil {
			err = errors.New(fmt.Sprintf("unable to locate sms no.: %v", err))
			return nil, err
		}
		_, err = RequestOtpRpc(&RequestOtpInput{Username: in.Username, Mobile: sendotp.Mobile, RequestType: in.RequestType})

		if err != nil {
			return nil, err
		}
		return &Jwtdata{Token: "proceed"}, nil

	}

	return nil, errors.New("invalid request type")

}
