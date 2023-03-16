package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
)



func (s *Server) SignUp(ctx context.Context, newUserData *NewUserDataInput) (*Jwtdata, error) { // takes id and sets up bucket and mongodb

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// search for duplicate username 
	//TODO change this to a map rather than search all docs
	result := MongoFields{}

	err := collection.FindOne(ctxMongo, bson.M{"Username": newUserData.Username}).Decode(&result)

	if err == nil {
		err = errors.New("username in use")
		return nil, err
	}

	createuser := MongoFields{Username: newUserData.Username, Email: newUserData.Email, Password: "depreciated", Profpic: "https://adminajh46unique.s3.eu-west-2.amazonaws.com/default-profile-pic.jpg", Photos: []string{}, LastCommentNum: 0, Posts: []*PostData{} }

	//username not in use so add new userdata struct
	_, err = collection.InsertOne(context.TODO(), createuser)
	if err != nil {
		return nil, err
	}

	passwordHash := utils.HashAndSalt([]byte(newUserData.Password))
	
	db := utils.Client.Database("datingapp").Collection("security")

	security :=	Security{Username: newUserData.Username, Password: passwordHash, OTP: OTP{}}

	_ , err = db.InsertOne(context.TODO(), security)
	if err != nil {
		return nil, err
	}


	utils.Createbucket(newUserData.Username) // create bucket to store users files

	//add error return when coial package gets pushed
	token, err2 := MakeJwt(&newUserData.Username, true) // make jwt with user id and auth true

	if err2 != nil {
		return nil, err2
	}

	return &Jwtdata{Token: token}, err2
}
