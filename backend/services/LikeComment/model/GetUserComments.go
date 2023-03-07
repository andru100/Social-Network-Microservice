package model

import (
	"context"
	//"fmt"
	//"net"
	"errors"
	//"log"
	//"net/http"
	"sort"
	"time"
	//"google.golang.org/grpc"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	//"github.com/andru100/Social-Network-Microservices/GetUserComments/model"
	//"github.com/andru100/Graphql-Social-Network/graph/model"
	"github.com/andru100/Social-Network/backend/social"
	//"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Server) GetUserComments(ctx context.Context, in *GetComments) (*MongoFields, error) {

	
	collection := social.Client.Database("datingapp").Collection("userdata") // connect to db and collection.
	currentDoc := MongoFields{}
	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&currentDoc)
	if err != nil {
		err = errors.New("unable to find users data")
		return nil, err
	}

	sort.Slice(currentDoc.Posts, func(i, j int) bool { // needs to be done on adding post and remove this
		return currentDoc.Posts[i].TimeStamp > currentDoc.Posts[j].TimeStamp
	})

	return &currentDoc, err
}
