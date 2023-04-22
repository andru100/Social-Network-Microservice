package model

import (
	"context"
	"errors"
	"sort"
	"fmt"
	"time"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/utils"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)


func GetReplys(ctx context.Context, in *GetPost) (*MongoFields, error) { // gets comments for all friends/users, for the home page feed

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.

	userdata := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var allPosts []*PostData

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&userdata)

	for _, data := range userdata.Replys {
		record := MongoFields{}
		err := collection.FindOne(ctxMongo, bson.M{"Username": data.Username}).Decode(&record)
		if err != nil {
			return nil, errors.New("unable to find users data")
		}
		for _, posts := range record.Posts {
			if posts.ID == data.PostID {
				allPosts = append(allPosts, posts)
			}
		}
	}

	//Sort posts or comments by time descending
	sort.Slice(allPosts, func(i, j int) bool { 
		return allPosts[i].TimeStamp > allPosts[j].TimeStamp
	})

	fmt.Println("sending users reply posts: ", allPosts)
	userdata.Posts = allPosts

	return &userdata, err
	

}
