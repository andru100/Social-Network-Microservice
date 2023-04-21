package model

import (
	"context"
	"fmt"
	"errors"
	"time"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserPosts(ctx context.Context, in *GetPost) (*MongoFields, error) {
	
	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.
	currentDoc := MongoFields{}
	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	if err != nil {
		err5 := errors.New("unable to find users data")
		fmt.Println(err5, err, in.Username)
		return nil, err5
	}
	
	return &currentDoc, err
}
