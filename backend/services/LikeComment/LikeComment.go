package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"errors"
	"time"
	"google.golang.org/grpc"
	"github.com/andru100/Social-Network-Microservices/backend/services/LikeComment/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/LikeComment/utils"
	"go.mongodb.org/mongo-driver/bson"
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

	collection := utils.Client.Database("datingapp").Collection("userdata")

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	
	//Find the comment being liked
	//TODO make each sub field a mongo doc so that comments can be searched by ID of doc and save on looping through data
	for i := 0; i < len(currentDoc.Posts); i++ {
		if currentDoc.Posts[i].ID == in.PostID {

			switch in.RequestType {
				case "create":
					likesent := model.Likes{
						Username: in.LikedBy ,
						Profpic:  in.LikeByPic,
					}
					currentDoc.Posts[i].Likes = append(currentDoc.Posts[i].Likes, &likesent) // add like to post

				case "delete":
					for j := 0; j < len(currentDoc.Posts[i].Likes); j++ {
						if currentDoc.Posts[i].Likes[j].Username == in.LikedBy {
							currentDoc.Posts[i].Likes = append(currentDoc.Posts[i].Likes[:j], currentDoc.Posts[i].Likes[j+1:]...) // delete like from post
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
				err = errors.New("error when adding Like to DB")
				return nil, err
			}
		}
	}

	go  LogUserLikes(in) // log like in user data

	
	// check originating page request came from and return updated comments to save an extra api call on react refresh component
	result, err1:= model.GetPostsClient (&model.GetPost{Username: in.Username, RequestType: in.ReturnPage})
	return result, err1
    
}

func LogUserLikes(in *model.SendLikeInput) error {
	collection := utils.Client.Database("datingapp").Collection("userdata")

	currentDoc := model.MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.LikedBy}).Decode(&currentDoc)
	
	
	switch in.RequestType {
		case "create":
			likesent := model.LikedData{
				Username: in.Username,
				PostID:   in.PostID,
			}

			currentDoc.Liked = append(currentDoc.Liked, &likesent) // add like to post

		case "delete":
			for j := 0; j < len(currentDoc.Liked); j++ {
				if currentDoc.Liked[j].PostID == in.PostID {
					currentDoc.Liked = append(currentDoc.Liked[:j], currentDoc.Liked[j+1:]...) // delete like from post
					break
				}
			}
		
		default:
			return  errors.New("invalid request type")
	}
			

	filter := bson.M{"Username": in.LikedBy}   
	Updatetype := "$set"
	Key2updt := "Liked"
	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, currentDoc.Liked},
		}},
	}

	//put to db
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.New("error when adding Like to DB")
	}

	return nil
		


	
}
