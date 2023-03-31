package model

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

func MakeJwt(userid *string, isauth bool) (string, error) {
	mySigningKey := []byte("AllYourBase") //base for ecoding eg private key

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(200 * time.Minute)
	
	// Create the JWT claims, which includes the username and expiry time
	claims := ClaimsChk{
		Username: *userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
	ss, err := token.SignedString(mySigningKey)                // encode to string with chosen base string

	return ss, err
}
