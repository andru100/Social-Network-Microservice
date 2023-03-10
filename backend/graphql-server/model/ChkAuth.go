package model

import (
	"fmt"
	"golang.org/x/net/context"
	jwt "github.com/dgrijalva/jwt-go"
)


func (s *Server) Chkauth(ctx context.Context, in *JwtdataInput) (*Authd, error) {
	
	fmt.Println("token iss", in.Token)

	var jwtKey = []byte("AllYourBase")

	claims := &ClaimsChk{}

	tkn, err := jwt.ParseWithClaims(*&in.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	var auth Authd

	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	auth.AuthdUser = claims.Username

	return &auth, err
}
