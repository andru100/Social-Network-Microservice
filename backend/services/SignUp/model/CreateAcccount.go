package model 

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
)

func CreateAccount (username string) error {
	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db

	createuser := MongoFields{Username: username, Profpic: "https://adminajh46unique.s3.eu-west-2.amazonaws.com/default-profile-pic.jpg", Photos: []string{}, LastCommentNum: 0, Posts: []*PostData{}}

	//username not in use so add new userdata struct
	_, err := collection.InsertOne(context.TODO(), createuser)
	if err != nil {
		return err
	}

	// security is passed so move from tempuser and add to security collection

	tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	temp2real := Security{}

	err = tempDB.FindOne(ctxMongo, bson.M{"Username": username}).Decode(&temp2real)

	if err != nil {
		return err
	}

	// expire otps

	temp2real.OTP.Mobile.Expiry = time.Now()
	temp2real.OTP.Email.Expiry = time.Now()

	db := utils.Client.Database("datingapp").Collection("security")

	_, err = db.InsertOne(context.TODO(), temp2real)
	if err != nil {
		return errors.New("its insertone on signup that is failing")
	}

	// delete from tempuser

	_, err = tempDB.DeleteOne(context.TODO(), bson.M{"Username": username})
	if err != nil {
		return errors.New("its deleteone on signup that is failing")
	}

	utils.Createbucket(username) // create bucket to store users files

	return err

}