package social

import (
	"fmt"
	"time"
	//"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/andru100/Graphql-Social-Network/graph/model"
)

func MakeJwt(userid *string, isauth bool) (string, error) {
	mySigningKey := []byte("AllYourBase") //base for ecoding eg private key

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	
	claims := model.Claims{
		Username: *userid,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // converts to json and signs/encodes
	ss, err := token.SignedString(mySigningKey)                // encode to string with chosen base string
	fmt.Printf("jwt created %v %v", ss, err)
	return ss, err

}
