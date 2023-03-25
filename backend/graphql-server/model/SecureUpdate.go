package model

import (
	"fmt"
	"errors"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
)

func (s *Server) SecureUpdate (ctx context.Context, in *SecurityCheckInput) (*Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	db := utils.Client.Database("datingapp").Collection("security")

	securityScore , err := SecurityCheck(in)

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

		token, err1 := MakeJwt(&in.Username, true)
		return &Jwtdata{Token: token}, err1

	} else {

		return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
	}
}
