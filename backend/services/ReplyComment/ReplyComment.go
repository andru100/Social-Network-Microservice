package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"errors"

	"google.golang.org/grpc"
	"github.com/google/uuid"
	"github.com/andru100/Social-Network-Microservices/backend/services/ReplyComment/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/ReplyComment/utils"
	"go.mongodb.org/mongo-driver/bson"
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

	fmt.Println("ReplyComment called!")
	
	uniqueID := uuid.New()
	
	collection := utils.Client.Database("datingapp").Collection("userdata")

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Find the document that mathes the id 
	err := collection.FindOne(ctxMongo, bson.M{"Username": in.AuthorUsername}).Decode(&currentDoc)

	//Find the post being replied to and add or delete it
	for i := 0; i < len(currentDoc.Posts); i++ {
		if currentDoc.Posts[i].ID == in.PostID {
			switch in.RequestType {
				case "create":
					

					reply := model.MsgCmts{
						ID:       uniqueID.String(),
						Username: in.ReplyUsername ,
						Comment:  in.ReplyComment , 
						Profpic:  in.ReplyProfpic ,
					}
					
					currentDoc.Posts[i].Comments = append(currentDoc.Posts[i].Comments, &reply) // add reply to post
				case "delete":
					for j := 0; j < len(currentDoc.Posts[i].Comments); j++ {
						if currentDoc.Posts[i].Comments[j].ID == in.ReplyID {
							currentDoc.Posts[i].Comments = append(currentDoc.Posts[i].Comments[:j], currentDoc.Posts[i].Comments[j+1:]...) // delete reply from post
							break
						}
					}
				default:
					err = errors.New("invalid request type")
					return nil, err
			}

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
			
			break
		}
	}

	go UpdateUserReplys(in, uniqueID.String()) // update replys to users record


	// check originating page request came from and return updated comments to save an extra api call on react refresh component
	username := in.ReplyUsername
	if in.ReturnPage == "user" {
		username = in.AuthorUsername
	} 
	result, err1:= model.GetPostsClient (&model.GetPost{Username: username, RequestType: in.ReturnPage})
	return result, err1
}

func UpdateUserReplys(in *model.ReplyCommentInput, uniqueId string) error {
	collection := utils.Client.Database("datingapp").Collection("userdata")

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Find the document that mathes the id 
	err := collection.FindOne(ctxMongo, bson.M{"Username": in.ReplyUsername}).Decode(&currentDoc)

	
	switch in.RequestType {
		case "create":
			replydata := model.ReplyData{
				Username: in.AuthorUsername,
				PostID: in.PostID,
				ReplyID: uniqueId, // this needs to be given after the reply is created and sent otherwise user commenting on self will screw it
			}

			currentDoc.Replys = append(currentDoc.Replys, &replydata) // add reply to post
		case "delete":
			for j := 0; j < len(currentDoc.Replys); j++ {
				if currentDoc.Replys[j].ReplyID == in.ReplyID {
					currentDoc.Replys = append(currentDoc.Replys[:j], currentDoc.Replys[j+1:]...) // delete reply from post
					break
				}
			}
		default:
			err = errors.New("invalid request type")
			return err
	}

	filter := bson.M{"Username": in.ReplyUsername}
	Updatetype := "$set"
	Key2updt := "Replys"
	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, currentDoc.Replys},
		}},
	}

	//put to db
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err = errors.New("error when adding comment to DB")
		return err
	} 

	return nil

	

}
