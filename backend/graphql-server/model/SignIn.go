package model

import (
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"golang.org/x/net/context"

)


func (s *Server) SignIn(ctx context.Context, in *SecurityCheckInput) (*Jwtdata, error) {// takes id and sets up bucket and mongodb doc

	// check username and password are correct and return security score. 

	switch in.RequestType {
		case "stage1":
			securityScore , err := SecurityCheck(in)
			// error will be throw if username or password is incorrect
			if err != nil {
				return nil, err
			}
			if securityScore >= 1 {
				result, err := RequestOtpRpc(&RequestOtpInput{Username: in.Username, Mobile: in.Mobile, RequestType: "sms"})

				if err != nil {
					return nil, err
				}
				return &Jwtdata{Token: "proceed"}, nil
			}
		case "stage2":

			securityScore , err := SecurityCheck(in)

			if securityScore >= 2 && err == nil {

				//generate jwt
				token, err1 := MakeJwt(&in.Username, true)
				return &Jwtdata{Token: token}, err1

			} else {

				return nil, errors.New(fmt.Sprintf("security check failed: %v", err))
			}
		default:
			return nil, errors.New("invalid stage")
	}
	return nil, errors.New("invalid stage")

}
