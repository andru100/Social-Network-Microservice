package main

import (
	"errors"
	"context"
	
	"golang.org/x/net/context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"

)

func (s *Server) SecureUpdate (ctx context.Context, in *SecurityCheckInput) (*Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	
	securityScore , err := utils.SecurityCheck(in)

	if securityScore >= 2 && err == nil {
		filter := bson.M{"Username": in.Username} 
 
		Updatetype := "$set"
		Key2updt := in.UpdateType
		update := bson.D{
			{Updatetype, bson.D{
				{Key2updt, result.OTP},
			}},
		}

		//put to db
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}

		token, err1 := MakeJwt(&in.Username, true)
		return &Jwtdata{Token: token}, err1

	} else {

		return nil, errors.New("security check failed: %v", err)
	}
}
