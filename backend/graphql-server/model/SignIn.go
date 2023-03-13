package model

import (
	"time"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
)


func (s *Server) SignIn(ctx context.Context, in *UsrsigninInput) (*Jwtdata, error) {

	db := utils.Client.Database("datingapp").Collection("security")

	result := Security{}

	ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

	err := db.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&result)

	if err != nil {
		return nil, errors.New("username not found")
	}

	fmt.Println("result", result)
	fmt.Println("password", in.Password)

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(in.Password))
    if err != nil {
        log.Println(err)
        return nil, errors.New("password does not match")
    } else {
		token, err1 := MakeJwt(&in.Username, true)
		return &Jwtdata{Token: token}, err1
	}


}
