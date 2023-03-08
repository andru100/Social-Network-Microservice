package model

import (
	"context"
	"log"
	"sort"
	"time"
	"github.com/andru100/Social-Network-Microservice/backend/GraphQL-Server/social"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Server) GetAllComments(ctx context.Context, in *GetComments) (*MongoFields, error) { // gets comments for all friends/users, for the home page feed

	collection := social.Client.Database("datingapp").Collection("userdata") // connect to db

	currentDoc := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var allPosts []*PostData

	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*MongoFields

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterating through cursor decod documents one at a time
	for cur.Next(context.TODO()) {

		var elem MongoFields
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	for _, record := range results {
		for _, posts := range record.Posts {
			allPosts = append(allPosts, posts)
		}
	}

	//Sort posts or comments by time descending
	//TODO sort posts as they are added to record to save on constantly resorting
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
