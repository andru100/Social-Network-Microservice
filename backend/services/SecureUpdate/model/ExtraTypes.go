package model 

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

type Server struct {
	UnimplementedSocialGrpcServer
}

type ClaimsChk struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}

type Security struct {
	Username string `json:"Username" bson:"Username"`
	Password Password `json:"Password" bson:"Password"`
	Email string `json:"Email" bson:"Email"`
	Mobile	string `json:"Mobile" bson:"Mobile"`
	DOB string `json:"DOB" bson:"DOB"`
	SecurityLock  SecurityLock `json:"SecurityLock" bson:"SecurityLock"`
	AuthType string `json:"AuthType" bson:"AuthType"`
	OTP 	OTP `json:"OTP" bson:"OTP"`
}	

type Password struct {
	Hash string `json:"Hash" bson:"Hash"`
	Attempts int `json:"Attempts" bson:"Attempts"`
}

type SecurityLock struct {
	Status string `json:"Status" bson:"Status"`
	Stage int `json:"Stage" bson:"Stage"`
	Expiry time.Time `json:"Expiry" bson:"Expiry"`
}

type OTP struct {
	Mobile Mobile `json:"Mobile" bson:"Mobile"`
	Email Email `json:"Email" bson:"Email"`
}	

type Mobile struct {
	Hash string `json:"OTP" bson:"OTP"`
	Expiry time.Time `json:"Expiry" bson:"Expiry"`
	Attempts int `json:"Attempts" bson:"Attempts"`
	Requests int `json:"Requests" bson:"Requests"`
}	

type Email struct {
	Hash string `json:"OTP" bson:"OTP"`
	Expiry time.Time `json:"Expiry" bson:"Expiry"`
	Attempts int `json:"Attempts" bson:"Attempts"`
	Requests int `json:"Requests" bson:"Requests"`
}	