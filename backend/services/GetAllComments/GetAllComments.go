package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"time"
	"google.golang.org/grpc"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network-Microservices/GetAllComments/model"
	"github.com/andru100/Social-Network/backend/social"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("GetAllComments running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4008))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	model.RegisterSocialGrpcServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (s *Server) GetAllComments(ctx context.Context, in *model.GetComments) (*model.MongoFields, error) { // gets comments for all friends/users, for the home page feed

	collection := social.Client.Database("datingapp").Collection("userdata") // connect to db and collection.

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var allPosts []*model.PostData

	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*model.MongoFields

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterating through cursor decode documents one at a time
	for cur.Next(context.TODO()) {

		var elem model.MongoFields
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

	var json2send model.MongoFields
	json2send.Posts = allPosts
	err = collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	json2send.Profpic = currentDoc.Profpic
	json2send.Bio = currentDoc.Bio
	json2send.Photos = currentDoc.Photos

	return &json2send, err
	

}
