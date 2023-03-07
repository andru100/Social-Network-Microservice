package model

import (
	"context"
	"time"
	"errors"
	"github.com/andru100/Social-Network-Microservice/social"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) NewComment (ctx context.Context, in *SendCmtInput) (*MongoFields, error) {
	
	collection := social.Client.Database("datingapp").Collection("userdata")

	currentDoc := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)

	currentDoc.LastCommentNum += 1
	
	//initialise empty slice to hold future likes and reply comments
	cmtHolder := []*MsgCmts{}
	likeHolder := []*Likes{}

	//make new comment struct: 
	newPost := PostData{
		Username:    in.Username,    
		SessionUser: in.SessionUser,
		MainCmt:     in.MainCmt,  
		PostNum:     currentDoc.LastCommentNum,    
		Time:        in.Time,   
		TimeStamp:   in.TimeStamp,    
		Date:        in.Date,    
		Comments:    cmtHolder ,
		Likes:      likeHolder, 
	}

	currentDoc.Posts = append(currentDoc.Posts, &newPost)
	filter := bson.M{"Username": in.Username} 

	//add new comment to DB 
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

	//update post index count
	update = bson.D{
		{Updatetype, bson.D{
			{"LastCommentNum", currentDoc.LastCommentNum},
		}},
	}

	//put to db
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err = errors.New("error when updating post index")
		return nil, err
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