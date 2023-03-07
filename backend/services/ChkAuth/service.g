
package model

import (
	"log"
	"golang.org/x/net/context"
)

type Server struct {
	UnimplementedChatServiceServer
}

func (s *Server) Chkauth(ctx context.Context, in *JwtdataInput) (*Authd, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Authd{Body: "Hello From the Server!"}, nil
}