package model

import (
	"context"
	"errors"
	"time"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)


func (s *Server) LikeComment (ctx context.Context, in *SendLikeInput) (*MongoFields, error) {

	collection := utils.Client.Database("datingapp").Collection("userdata")

	currentDoc := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	
	//Find the comment being liked
	//TODO make each sub field a mongo doc so that comments can be searched by ID of doc and save on looping through data
	for i := 0; i < len(currentDoc.Posts); i++ {
		if currentDoc.Posts[i].ID == in.PostID {
			likesent := Likes{
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
		result, err1 := Rpc2GetAllCmts(&GetComments{Username: in.Username})
		return result, err1
    } else {
	   result, err1:= Rpc2GetUserCmts (&GetComments{Username: in.Username})
	   return result, err1
    }
}
