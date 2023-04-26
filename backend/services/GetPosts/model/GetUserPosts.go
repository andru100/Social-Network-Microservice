package model

import (
	"context"
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
		return nil, errors.New("unable to find users data")
	}
	
	return &currentDoc, err
}
