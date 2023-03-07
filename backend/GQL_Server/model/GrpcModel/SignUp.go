package model

import (
	"context"
	"errors"
	"time"
	"github.com/andru100/Social-Network-Microservice/backend/GQL_Server/social"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) SignUp(ctx context.Context, newUserData *NewUserDataInput) (*Jwtdata, error) { // takes id and sets up s3 bucket and db

	collection := social.Client.Database("datingapp").Collection("userdata") // connect to db

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// search for duplicate username 
	//TODO change this to a map rather than search all docs
	result := MongoFields{}

	err := collection.FindOne(ctxMongo, bson.M{"Username": newUserData.Username}).Decode(&result)

	if err == nil {
		err = errors.New("username in use")
		return nil, err
	}

	createuser := Usrsignin{Username: newUserData.Username, Email: newUserData.Email, Password: newUserData.Password, Photos: []string{}, LastCommentNum: 0, Posts: []*PostData{} }

	//username not in use so add new userdata struct
	_, err = collection.InsertOne(context.TODO(), createuser)
	if err != nil {
		err = errors.New("problem creating user")
		return nil, err
	}


	social.Createbucket(newUserData.Username) // create bucket to store users files

	//add error return when social package gets pushed
	token := social.MakeJwt(&newUserData.Username, true) // make jwt with user id and auth true

	if err != nil {
		return nil, err
	}

	return &Jwtdata{Token: token}, err
}
