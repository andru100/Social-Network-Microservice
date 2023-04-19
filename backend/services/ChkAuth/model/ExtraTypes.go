package model 

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type ClaimsChk struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}