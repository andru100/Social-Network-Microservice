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
	Password string `json:"Password" bson:"Password"`
	OTP 	OTP `json:"OTP" bson:"OTP"`
}	

type OTP struct {
	Mobile Mobile `json:"Mobile" bson:"Mobile"`
	Email Email `json:"Email" bson:"Email"`
}	

type Mobile struct {
	Hash string `json:"OTP" bson:"OTP"`
	Expiry time.Time `json:"Expiry" bson:"Expiry"`
	Attempts int `json:"Attempts" bson:"Attempts"`
}	

type Email struct {
	Hash string `json:"OTP" bson:"OTP"`
	Expiry time.Time `json:"Expiry" bson:"Expiry"`
	Attempts int `json:"Attempts" bson:"Attempts"`
}	

