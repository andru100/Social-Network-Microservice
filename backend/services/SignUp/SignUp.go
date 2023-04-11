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
	fmt.Println("signup called request is", in)
	
	switch in.RequestType {
	case "username":
		fmt.Println("signup username called:", in, "in.UpdateData is: ", in.UpdateData)
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

		if err != nil {
			tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

			ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

			// search for duplicate username
			//TODO change this to a map rather than search all docs
			verifyTempUsername := model.Security{}
			err = tempDB.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&verifyTempUsername)

			if err == nil {
				err = errors.New("someone is already in the signup process with this username")
				return nil, err
			}

			passwordHash := utils.HashAndSalt([]byte(in.Password))

			passwordHolder := model.Password{Hash: passwordHash, Attempts: 0}

			tempuser := model.Security{Username: in.Username, Password: passwordHolder, DOB: in.DOB, SecurityLevel: 100 }

			_, err = tempDB.InsertOne(context.TODO(), tempuser)
			if err != nil {
				return nil, errors.New("its insertone on signup that is failing")
			}

			return &model.Jwtdata{Token: "proceed", AuthType: "email"}, nil
		}
	case "email":
		

		db := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)
		fmt.Println("signup email called:", in, "in.UpdateData is: ", in.UpdateData)
		verifyEmail := model.Security{}

		err := db.FindOne(ctxMongo, bson.M{"Email": in.Email}).Decode(&verifyEmail)

		if err == nil {
			err = errors.New("email in use")
			return nil, err
		} else {
			tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.
			filter := bson.M{"Username": in.Username} 
		
		
			update := bson.D{
				{"$set", bson.D{
					{"Email", in.Email},
				}},
			}

			_, err = tempDB.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, RequestType: "email", UserType: "temp"})

			if err != nil {
				fmt.Println(err)
				return nil, errors.New("error requesting email otp")
			}

			//mobileclue := tempuser.Mobile[len(tempuser.Mobile)-3:] 
			emailclue := in.Email[0:3]

			return &model.Jwtdata{Token: "proceed", AuthType: "confirmemail", EmailClue: emailclue}, nil
				
		}
	case "sms":
		fmt.Println("signup mobile called:", in, "in.UpdateData is: ", in.UpdateData)
		verifyMobile := model.Security{}

		db := utils.Client.Database("datingapp").Collection("security") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

		err := db.FindOne(ctxMongo, bson.M{"Mobile": in.Mobile}).Decode(&verifyMobile)

		if err == nil {
			fmt.Println("error found dupliate mobile in temp signup", verifyMobile)
			err = errors.New("mobile in use")
			return nil, err
		} else {
			tempdb := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.
			filter := bson.M{"Username": in.Username} 
		
		
			update := bson.D{
				{"$set", bson.D{
					{"Mobile", in.Mobile},
				}},
			}

			_, err = tempdb.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, RequestType: "sms", UserType: "temp"})

			if err != nil {
				fmt.Println(err)
				return nil, errors.New("error requesting sms otp")
			}

			//mobileclue := tempuser.Mobile[len(tempuser.Mobile)-3:] 
			emailclue := in.Email[0:3]
			mobclue := in.Mobile[len(in.Mobile)-3:]

			return &model.Jwtdata{Token: "proceed", AuthType: "confirmsms", EmailClue: emailclue, MobClue: mobclue}, nil	
		}
	case "setsecurity":
		fmt.Println("setsecurity called:", in, )
		tempDB := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)
		verifySecurity := model.Security{}

		err := tempDB.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&verifySecurity)

		if err != nil {
			fmt.Println("error finding user in setsecurity")
			return nil, err
		} 

		if verifySecurity.SecurityLevel == 100 { // if is unset then set security level

		
			securitymatrix := make(map[string]int)
			securitymatrix["email"] = 2
			securitymatrix["sms"] = 2
			securitymatrix["password"] = 1
			securitymatrix["oauth"] = 0
			securitymatrix["high"] = 3

			userlevel := securitymatrix[in.Token]
			

			//verifySecurity.SecurityLevel = securitymatrix[in.UpdateData] bencmmark updating in place 

			filter := bson.M{"Username": in.Username} 

			update := bson.D{
				{"$set", bson.D{
					{"SecurityLevel", userlevel},
				}},
			}

			_, err = tempDB.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}

			update = bson.D{
				{"$set", bson.D{
					{"AuthType", in.Token},
				}},
			}

			_, err = tempDB.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, err
			}
			
			securityscore, err := model.SecurityCheck(in)
			if err != nil {
				return nil, err
			}
			if securityscore >= userlevel+1 || in.Token == "email" && securityscore >= userlevel { //+1 because includes password wwill remove from react state.. allow for email becauuse dnt need a +! all other auth types collect email aswell

				err = model.CreateAccount(in.Username)
				if err != nil {
					errors.New("error creating account")
					return nil, err
				}

				token, err2 := model.MakeJwt(&in.Username, true) // make jwt with user id and auth true
				fmt.Println("made jwt in ccreat acccount")

				if err2 != nil {
					return nil, err2
				}

				return &model.Jwtdata{Token: token}, err2
			}

			return &model.Jwtdata{Token: "proceed", AuthType: in.Token}, nil
		}

	

	case "stage2":

		fmt.Println("signup stage 2 new called", in)

		collection := utils.Client.Database("datingapp").Collection("tempuser") // connect to db and collection.

		ctxMongo, _ := context.WithTimeout(context.Background(), 15*time.Second)

		// search for duplicate username
		//TODO change this to a map rather than search all docs
		userdata := model.Security{}

		err := collection.FindOne(ctxMongo, bson.M{"Username": in.Username}).Decode(&userdata)

		if err != nil {
			return nil, err
		}
		
		
		securityscore, err := model.SecurityCheck(in)

		if err != nil  || securityscore < 1{ //was 2
			return nil, errors.New(fmt.Sprintf("security check failed, score: %v error %v", securityscore, err))
		}

		if securityscore >1 && userdata.SecurityLevel == 100 {
			return &model.Jwtdata{Token: "proceed", AuthType: "setsecurity"}, nil
		}

		

		if securityscore > 0 && securityscore < userdata.SecurityLevel { // if is  check request gets sccore 1 for pass and doesnt return errors then do a rpc and return a proceed and auth type being asked for

			fmt.Println("security check passed but not enough, sending next steps and otp")

			_, err = model.RequestOtpRpc(&model.RequestOtpInput{Username: in.Username, Email: in.Email, Mobile: in.Mobile, RequestType: userdata.AuthType, UserType: "temp"})

			if err != nil {
				fmt.Println(err)
				return nil, errors.New("error requesting otp on signup stage 2")
			}
			
			mobileclue := userdata.Mobile[len(userdata.Mobile)-3:] 
			emailclue := userdata.Email[0:3]
		
			return &model.Jwtdata{Token: "proceed", AuthType: userdata.AuthType, MobClue: mobileclue, EmailClue: emailclue}, nil
		}

		// add checcks in react for if token is proceed if not take it and go home. also change stages being called to all stage2


		if securityscore >= userdata.SecurityLevel  { // was 3 create user return token
			fmt.Println("security score is", securityscore, "user security level is", userdata.SecurityLevel)
			fmt.Println("verification passed, creating user")

			err = model.CreateAccount(in.Username)
			if err != nil {
				errors.New("error creating account")
				return nil, err
			}

			
			//add error return when coial package gets pushed
			token, err := model.MakeJwt(&in.Username, true) // make jwt with user id and auth true

			if err != nil {
				return nil, err
			}

			return &model.Jwtdata{Token: token}, err
		} else {
			return nil, errors.New(fmt.Sprintf("security check failed, score: %v ", securityscore))
		}

	default:
		return nil, errors.New("invalid request type")
	}
	return nil, errors.New("hate this   ")
}




