package model

import (
	"context"
	"time"
	"errors"

	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) ReplyComment (ctx context.Context, in *ReplyCommentInput) (*MongoFields, error) {// adds replies to users comments mongo doc

	collection := utils.Client.Database("datingapp").Collection("userdata")

	currentDoc := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Find the document
	err := collection.FindOne(ctxMongo, bson.M{"Username": in.AuthorUsername}).Decode(&currentDoc)

	//Find the comment being replied do by index and add it
	//TODO make each sub filed a mongo doc so we can search by ID and save looping through all comments
	for i := 0; i < len(currentDoc.Posts); i++ {
		if currentDoc.Posts[i].ID == in.PostID {
			reply := MsgCmts{
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
	result, err1:= GetPostsClient (&GetPost{Username: in.AuthorUsername, RequestType: in.ReturnPage})
	return result, err1
}
