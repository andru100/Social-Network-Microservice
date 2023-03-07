package social

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/andru100/Social-Network-Microservice/model"
)

func MakeJwt(userid *string, isauth bool) (string, error) {
	mySigningKey := []byte("AllYourBase") //base for ecoding eg private key

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	
	// Create the JWT claims, which includes the username and expiry time
	claims := model.ClaimsChk{
		Username: *userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
	ss, err := token.SignedString(mySigningKey)                // encode to string with chosen base string

	return ss, err
}
