package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"

	//"golang.org/x/crypto/bcrypt"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/model"
	"github.com/andru100/Social-Network-Microservices/backend/services/SignUp/utils"
)

type Server struct {
	model.UnimplementedSocialGrpcServer
}

func main() {

	fmt.Println("SignUp running!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4002))
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

func (s *Server) SignUp(ctx context.Context, in *model.SecurityCheckInput) (*model.Jwtdata, error) { // takes id and sets up bucket and mongodb
	fmt.Println("signup called request type is", in.RequestType)
	switch in.RequestType {
	case "stage1":
		collection := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

		// search for duplicate username
		//TODO change this to a map rather than search all docs
		verifyUsername := model.Security{}

		err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&verifyUsername)

		if err == nil {
			err = errors.New("username in use")
			return nil, err
		}

		// search for duplicate email

		verifyEmail := model.Security{}

		err = collection.FindOne(ctxMongo, bson.M{"Email": in.Email}).Decode(&verifyEmail)

		if err == nil {
			err = errors.New("email in use")
			return nil, err
		}

		verifyMobile := model.Security{}

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
		verifyTempUsername := model.Security{}

		err = tempDB.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&verifyTempUsername)

		if err == nil {
			err = errors.New("someone is already in the signup process with this username")
			return nil, err
		}

		// add other checks for email and mobile and work out logic

		// add to tempuser collection

		passwordHash := utils.HashAndSalt([]byte(in.Password))

		passwordHolder := model.Password{Hash: passwordHash, Attempts: 0}

		tempUser := model.Security{Username: in.Username, Password: passwordHolder, DOB: in.DOB, Email: in.Email, Mobile: in.Mobile, OTP: model.OTP{}}

		_, err = tempDB.InsertOne(context.TODO(), tempUser)
		if err != nil {
			return nil, errors.New("its insertone on signup that is failing")
		}

		// send otps

		_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, Mobile: in.Mobile, RequestType: "signup"})

		if err != nil {
			fmt.Println(err)
			return nil, errors.New("its requestotp on signup that is failing")
		}

		return &model.Jwtdata{Token: "proceed"}, err

	case "stage2":

		fmt.Println(*in)

		securityScore, err := model.SecurityCheck(in)

		fmt.Println("security score is", securityScore)


		if err != nil  || securityScore < 2 {
			return nil, errors.New(fmt.Sprintf("security check failed, score: %v error %v", securityScore, err))
		}

		if securityScore >= 2  {
			return &model.Jwtdata{Token: "proceed"}, err
		}


	case "stage3":

		fmt.Println(in)

		securityScore, err := model.SecurityCheck(in)

		if err != nil  || securityScore < 2 {
			return nil, errors.New(fmt.Sprintf("security check failed, score: %v error %v", securityScore, err))
		}

		if securityScore >= 3 {

			collection := utils.Client.Database("datingapp").Collection("userdata") // connect to db

			createuser := model.MongoFields{Username: in.Username, Profpic: "https://adminajh46unique.s3.eu-west-2.amazonaws.com/default-profile-pic.jpg", Photos: []string{}, LastCommentNum: 0, Posts: []*model.PostData{}}

			//username not in use so add new userdata struct
			_, err = collection.InsertOne(context.TODO(), createuser)
			if err != nil {
				return nil, err
			}

			// security is passed so move from tempuser and add to security collection

			tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

			ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

			temp2real := model.Security{}

			err = tempDB.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&temp2real)

			if err != nil {
				return nil, err
			}

			// expire otps

			temp2real.OTP.Mobile.Expiry = time.Now()
			temp2real.OTP.Email.Expiry = time.Now()

			db := utils.Client.Database("datingapp").Collection("security")

			_, err = db.InsertOne(context.TODO(), temp2real)
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
			token, err2 := model.MakeJwt(&in.Username, true) // make jwt with user id and auth true

			if err2 != nil {
				return nil, err2
			}

			return &model.Jwtdata{Token: token}, err2
		} else {
			return nil, err
		}

	default:
		return nil, errors.New("invalid request type")
	}
	return nil, errors.New("hate this   ")
}
