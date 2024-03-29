package model

import (
	"context"
	"errors"
	"sort"
	"time"
	
	"github.com/andru100/Social-Network-Microservices/backend/services/ReplyComment/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) GetPosts(ctx context.Context, in *GetPost) (*MongoFields, error) {

	
	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.
	currentDoc := MongoFields{}
	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	if err != nil {
		err = errors.New("unable to find users data")
		return nil, err
	}

	sort.Slice(currentDoc.Posts, func(i, j int) bool { // needs to be done on adding post and remove this
		return currentDoc.Posts[i].TimeStamp > currentDoc.Posts[j].TimeStamp
	})

	return &currentDoc, err
}
