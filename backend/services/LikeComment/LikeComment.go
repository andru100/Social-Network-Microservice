package main

import (
	"context"
	"fmt"
	"log"
	"net"
	// "sort"
	"errors"
	"time"
	"google.golang.org/grpc"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network-Microservices/LikeComment/model"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network/backend/social"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("LikeComment running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4003))
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


func (s *Server) LikeComment (ctx context.Context, in *model.SendLikeInput) (*model.MongoFields, error) {

	collection := social.Client.Database("datingapp").Collection("userdata")

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	
	//Find the comment being liked
	//TODO make each sub field a mongo doc so that comments can be searched by ID of doc and save on looping through data
	for i := 0; i < len(currentDoc.Posts); i++ {
		if currentDoc.Posts[i].PostNum == in.PostIndx {
			likesent := model.Likes{
				Username: in.LikedBy ,
				Profpic:  in.LikeByPic,
			}
			currentDoc.Posts[i].Likes = append(currentDoc.Posts[i].Likes, &likesent) // add like to post
			filter := bson.M{"Username": currentDoc.Posts[i].Username}   
			Updatetype := "$set"
			Key2updt := "Posts"
			update := bson.D{
				{Updatetype, bson.D{
					{Key2updt, currentDoc.Posts},
				}},
			}
		
			//put to db
			_, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				err = errors.New("error when adding Like to DB")
				return nil, err
			}
		}
	}

	
	// check originating page request came from and return updated comments to save an extra api call on react refresh component
	if in.ReturnPage == "All" {
		result, err1 := model.Rpc2GetAllCmts(&model.GetComments{Username: in.Username})
		return result, err1
    } else {
	   result, err1:= model.Rpc2GetUserCmts (&model.GetComments{Username: in.Username})
	   return result, err1
    }
}
