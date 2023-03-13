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
	OTP string `json:"OTP" bson:"OTP"`
	OTPExpire time.Time `json:"OTPExpire" bson:"OTPExpire"`
}	