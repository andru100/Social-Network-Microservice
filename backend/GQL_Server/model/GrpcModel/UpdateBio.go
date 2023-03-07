package model

import (
	"context"
	"errors"
	"time"
	"github.com/andru100/Social-Network-Microservice/social"
	"go.mongodb.org/mongo-driver/bson"
)


func (s *Server) UpdateBio (ctx context.Context, in *UpdateBioInput) (*MongoFields, error) { // updates user bio section
	
	collection := social.Client.Database("datingapp").Collection("userdata")

	filter := bson.M{"Username": in.Username}

	Updatetype := "$set"
	Key2updt := "Bio"

	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, in.Bio},
		}},
	}

	//put to db
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err = errors.New("error when updating to DB")
		return nil, err
	}

	currentDoc := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err = collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)

	return &currentDoc, err
}
