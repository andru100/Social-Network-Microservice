package model

import (
	"context"
	"errors"
	"time"
	"fmt"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)


func (s *Server) Follow (ctx context.Context, in *FollowInput) (*MongoFields, error) {
	fmt.Println("Follow/unfollow request recieved")

	collection := utils.Client.Database("datingapp").Collection("userdata")

	requester := MongoFields{}
	user := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&requester)

	if err != nil {
		err = errors.New("unable to find users data")
		return nil, err
	}

	err = collection.FindOne(ctxMongo, bson.M{"Username": in.UserOfIntrest}).Decode(&user)

	if err != nil {
		err = errors.New("unable to find users data")
		return nil, err
	}

	switch in.RequestType {
		case "follow":
			requester.Following = append(requester.Following, in.UserOfIntrest) 
			user.Followers = append(user.Followers, in.Username)

		case "unfollow":
			for i := 0; i < len(requester.Following); i++ {
				if requester.Following[i] == in.UserOfIntrest {
					requester.Following = append(requester.Following[:i], requester.Following[i+1:]...)
					break
				}
			}

			for i := 0; i < len(user.Followers); i++ {
				if user.Followers[i] == in.Username {
					user.Followers = append(user.Followers[:i], user.Followers[i+1:]...)
					break
				}
			}
		
		default:
			err = errors.New("invalid request type")
			return nil, err
	}

	filter := bson.M{"Username": in.Username}
	Updatetype := "$set"
	Key2updt := "Following"
	update := bson.D{
		{Updatetype, bson.D{
			{Key2updt, requester.Following},
		}},
	}
		
	//put to db
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err = errors.New("error when adding Like to DB")
		return nil, err
	}

	filter = bson.M{"Username": in.UserOfIntrest}
	Updatetype = "$set"
	Key2updt = "Followers"
	update = bson.D{
		{Updatetype, bson.D{
			{Key2updt, user.Followers},
		}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		err = errors.New("error when adding Like to DB")
		return nil, err
	}

		
	// check originating page request came from and return updated comments to save an extra api call on react refresh component
	result, err1:= GetPostsClient (&GetPost{Username: in.Username, RequestType: in.ReturnPage})
	return result, err1
}


