package model

import (
	"context"
	"errors"
	"time"
	"log"
	"golang.org/x/crypto/bcrypt"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) SignUp(ctx context.Context, newUserData *NewUserDataInput) (*Jwtdata, error) { // takes id and sets up s3 bucket and db

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// search for duplicate username 
	//TODO change this to a map rather than search all docs
	result := MongoFields{}

	err := collection.FindOne(ctxMongo, bson.M{"Username": newUserData.Username}).Decode(&result)

	if err == nil {
		err = errors.New("username in use")
		return nil, err
	}

	createuser := Usrsignin{Username: newUserData.Username, Email: newUserData.Email, Password: "depreciated", Photos: []string{}, LastCommentNum: 0, Posts: []*PostData{} }

	//username not in use so add new userdata struct
	_, err = collection.InsertOne(context.TODO(), createuser)
	if err != nil {
		err = errors.New("problem creating user")
		return nil, err
	}

	password := hashAndSalt([]byte(newUserData.Password))
	
	db := utils.Client.Database("datingapp").Collection("security")

	security :=	Security{Username: newUserData.Username, Password: password, OTP: OTP{}}

	_ , err = db.InsertOne(context.TODO(), security)
	if err != nil {
		return nil, err
	}


	utils.Createbucket(newUserData.Username) 

	token, err1 := MakeJwt(&newUserData.Username, true)
	if err1 != nil {
		return nil, err1
	}

	return &Jwtdata{Token: token}, err1
}


func hashAndSalt(pwd []byte) string {
    
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }
    return string(hash)
}
