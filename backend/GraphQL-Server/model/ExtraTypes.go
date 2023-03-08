package model 

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Server struct {
	UnimplementedSocialGrpcServer
}

type ClaimsChk struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}