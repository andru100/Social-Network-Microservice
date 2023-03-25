package model

import (
	"context"
	"errors"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	//"golang.org/x/crypto/bcrypt"
	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
)


func (s *Server) SignUp(ctx context.Context, in *SecurityCheckInput) (*Jwtdata, error) { // takes id and sets up bucket and mongodb

	switch in.RequestType {
		case "stage1":
			collection := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

			ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

			// search for duplicate username 
			//TODO change this to a map rather than search all docs
			verifyUsername := Security{}

			err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&verifyUsername)

			if err == nil {
				err = errors.New("username in use")
				return nil, err
			}

			// search for duplicate email

			verifyEmail := Security{}

			err = collection.FindOne(ctxMongo, bson.M{"Email": in.Email}).Decode(&verifyEmail)

			if err == nil {
				err = errors.New("email in use")
				return nil, err
			}

			verifyMobile := Security{}

			err = collection.FindOne(ctxMongo, bson.M{"Mobile": in.Email}).Decode(&verifyMobile)

			if err == nil {
				err = errors.New("mobile in use")
				return nil, err
			}

			// no duplicate found so ping requstotp flag as temp add to tempuser collection 

			tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

			ctxMongo, _ = context.WithTimeout(context.Background(), 15*time.Second)

			// search for duplicate username 
			//TODO change this to a map rather than search all docs
			verifyTempUsername := Security{}

			err = tempDB.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&verifyTempUsername)

			if err == nil {
				err = errors.New("someone is already in the signup process with this username")
				return nil, err
			}

			// add other checks for email and mobile and work out logic

			// add to tempuser collection

			passwordHash := utils.HashAndSalt([]byte(in.Password))

			passwordHolder := Password{Hash: passwordHash, Attempts: 0}

			tempUser := Security{Username: in.Username, Password: passwordHolder, DOB: in.DOB, Email: in.Email, Mobile: in.Mobile, OTP: OTP{}}

			_ , err = tempDB.InsertOne(context.TODO(), tempUser)
			if err != nil {
				return nil, errors.New("its insertone on signup that is failing")
			}

			// send otps

			_, err = RequestOtpRpc(&RequestOtpInput{Username: in.Username, Email: in.Email, Mobile: in.Mobile, RequestType: "signup"})

			if err != nil {
				fmt.Println(err)
				return nil, errors.New("its requestotp on signup that is failing")
			}

			return &Jwtdata{Token: "proceed"}, err





		case "stage2":

			securityScore , err := SecurityCheck(in)

			if securityScore >= 2 && err == nil {

				collection := utils.Client.Database("datingapp").Collection("users") // connect to db
				
				createuser := MongoFields{Username: in.Username, Profpic: "https://adminajh46unique.s3.eu-west-2.amazonaws.com/default-profile-pic.jpg", Photos: []string{}, LastCommentNum: 0, Posts: []*PostData{} }

				//username not in use so add new userdata struct
				_, err = collection.InsertOne(context.TODO(), createuser)
				if err != nil {
					return nil, err
				}

				// security is passed so move from tempuser and add to security collection

				tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

				ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

				temp2real := Security{}

				err = tempDB.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&temp2real)

				if err == nil {
					err = errors.New("if this happens youve been hacked")
					return nil, err
				}
				
				db := utils.Client.Database("datingapp").Collection("security")

				_ , err = db.InsertOne(context.TODO(), temp2real)
				if err != nil {
					return nil, errors.New("its insertone on signup that is failing")
				}

				// delete from tempuser

				_, err = tempDB.DeleteOne(context.TODO(), bson.M{"Username": in.Username})
				if err != nil {
					return nil, errors.New("its deleteone on signup that is failing")
				}


				utils.Createbucket(in.Username) // create bucket to store users files

				//add error return when coial package gets pushed
				token, err2 := MakeJwt(&in.Username, true) // make jwt with user id and auth true

				if err2 != nil {
					return nil, err2
				}

				return &Jwtdata{Token: token}, err2
			} else {
				return nil, err
			}
			
		default:
			return nil, errors.New("invalid request type")
	}
}
