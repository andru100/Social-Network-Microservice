package main

import (
	"context"
	"fmt"
	"log"
	"net"
	//"sort"
	"time"
	"errors"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/ReplyComment/model"
	"github.com/andru100/Social-Network/backend/social"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("ReplyComment running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4004))
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


func (s *Server) ReplyComment (ctx context.Context, in *model.ReplyCommentInput) (*model.MongoFields, error) {// adds replies to users comments mongo doc

	collection := social.Client.Database("datingapp").Collection("userdata")

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Find the document that mathes the id 
	err := collection.FindOne(ctxMongo, bson.M{"Username": in.AuthorUsername}).Decode(&currentDoc)

	//Find the comment being replied do by index and add it
	//TODO make each sub filed a mongo doc so we can search by ID and save looping through all comments
	for i := 0; i < len(currentDoc.Posts); i++ {
		if currentDoc.Posts[i].PostNum == in.PostIndx {
			reply := model.MsgCmts{
				Username: in.ReplyUsername ,
				Comment:  in.ReplyComment , 
				Profpic:  in.ReplyProfpic ,
			}
			currentDoc.Posts[i].Comments = append(currentDoc.Posts[i].Comments, &reply) // add reply to post
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
				err = errors.New("error when adding comment to DB")
				return nil, err
			}                             
		}
	}

	
	// check originating page request came from and return updated comments to save an extra api call on react refresh component
	if in.ReturnPage == "All" {
		result, err1 := model.Rpc2GetAllCmts(&model.GetComments{Username: in.ReplyUsername})
		return result, err1
    } else {
	   result, err1:= model.Rpc2GetUserCmts (&model.GetComments{Username: in.AuthorUsername})
	   return result, err1
    }
}
