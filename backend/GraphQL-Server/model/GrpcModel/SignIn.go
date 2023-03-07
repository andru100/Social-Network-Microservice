package model

import (
	"time"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"github.com/andru100/Social-Network-Microservice/backend/GraphQL-Server/social"
)


func (s *Server) SignIn(ctx context.Context, in *UsrsigninInput) (*Jwtdata, error) {

	collection := social.Client.Database("datingapp").Collection("userdata") // connect to db

	result := MongoFields{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	fmt.Println("result", result)
	fmt.Println("password", in.Password)
	
	if err != nil {
		return nil, errors.New("username not found")
	}

	if result.Password == in.Password {
		token := social.MakeJwt(&in.Username, true)
		return &Jwtdata{Token: token}, nil
	} else {
		return nil, errors.New("password does not match")
	}

}
