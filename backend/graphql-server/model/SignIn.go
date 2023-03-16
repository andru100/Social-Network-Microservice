package model

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	//"github.com/andru100/Social-Network-Microservice/backend/graphql-server/utils"
)



func (s *Server) SignIn(ctx context.Context, in *SecurityCheckInput) (*Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	securityScore , err := SecurityCheck(in)

	if securityScore >= 2 && err == nil {

		//generate jwt
		token, err1 := MakeJwt(&in.Username, true)
		return &Jwtdata{Token: token}, err1

	} else {

		return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
	}

}
