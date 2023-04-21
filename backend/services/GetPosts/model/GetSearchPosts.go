package model

import (
	"context"
	"log"
	"sort"
	"time"
	"strings"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/GetUserComments/utils"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetSearchPosts(ctx context.Context, in *GetPost) (*MongoFields, error) { // gets comments for all friends/users, for the home page feed

	collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db and collection.

	currentDoc := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var allPosts []*PostData

	findOptions := options.Find()
	findOptions.SetLimit(2)

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterating through cursor decode documents one at a time
	for cur.Next(context.TODO()) {

		var userdata MongoFields
		err := cur.Decode(&userdata)
		if err != nil {
			log.Fatal(err)
		}

		for _, post := range userdata.Posts {
			if strings.Contains(post.MainCmt, in.SearchTerm) {
				allPosts = append(allPosts, post)
			}
		}

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	//Sort posts or comments by time descending
	sort.Slice(allPosts, func(i, j int) bool { 
		return allPosts[i].TimeStamp > allPosts[j].TimeStamp
	})

	var json2send MongoFields
	json2send.Posts = allPosts
	err = collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	json2send.Profpic = currentDoc.Profpic
	json2send.Bio = currentDoc.Bio
	json2send.Photos = currentDoc.Photos

	return &json2send, err
	

}
