package model

import (
	"context"
	"errors"
	"sort"
	"time"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/utils"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)


func GetFollowingPosts(ctx context.Context, in *GetPost) (*MongoFields, error) { // gets comments for all friends/users, for the home page feed

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.

	userdata := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var allPosts []*PostData

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&userdata)

	for _, user := range userdata.Following {
		record := MongoFields{}
		err := collection.FindOne(ctxMongo, bson.M{"Username": user}).Decode(&record)
		if err != nil {
			return nil, errors.New("unable to find users data")
		}
		for _, posts := range record.Posts {
			allPosts = append(allPosts, posts)
		}
	}

	//Sort posts or comments by time descending
	sort.Slice(allPosts, func(i, j int) bool { 
		return allPosts[i].TimeStamp > allPosts[j].TimeStamp
	})

	var json2send MongoFields
	json2send.Posts = allPosts
	json2send.Profpic = userdata.Profpic
	json2send.Bio = userdata.Bio
	json2send.Photos = userdata.Photos

	return &json2send, err
	

}
