package graph

import (
	"context"
	"log"
	"os"

	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/model"
	"google.golang.org/grpc"
)

// Signin is the resolver for the Signin field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.UsrsigninInput) (*model.Jwtdata, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4001", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.SignIn(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// SignUp is the resolver for the SignUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input model.NewUserDataInput) (*model.Jwtdata, error) {
	
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4002", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.SignUp(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// LikeComment is the resolver for the LikeComment field.
func (r *mutationResolver) LikeComment(ctx context.Context, input model.SendLikeInput) (*model.MongoFields, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4003", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.LikeComment(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// ReplyComment is the resolver for the ReplyComment field.
func (r *mutationResolver) ReplyComment(ctx context.Context, input model.ReplyCommentInput) (*model.MongoFields, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4004", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.ReplyComment(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// NewComment is the resolver for the NewComment field.
func (r *mutationResolver) NewComment(ctx context.Context, input model.SendCmtInput) (*model.MongoFields, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4005", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.NewComment(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// UpdateBio is the resolver for the UpdateBio field.
func (r *mutationResolver) UpdateBio(ctx context.Context, input model.UpdateBioInput) (*model.MongoFields, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4006", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.UpdateBio(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// Chkauth is the resolver for the Chkauth field.
func (r *queryResolver) Chkauth(ctx context.Context, input model.JwtdataInput) (*model.Authd, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4007", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.Chkauth(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result, err
}

// GetAllComments is the resolver for the GetAllComments field.
func (r *queryResolver) GetAllComments(ctx context.Context, username string) (*model.MongoFields, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4008", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.GetAllComments(context.Background(), &model.GetComments{Username: username})

	if err != nil {
		return nil, err
	}

	return result, err

}

// GetUserComments is the resolver for the GetUserComments field.
func (r *queryResolver) GetUserComments(ctx context.Context, username string) (*model.MongoFields, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(os.Getenv("HOSTIP")+":4009", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := model.NewSocialGrpcClient(conn)

	result, err := c.GetUserComments(context.Background(), &model.GetComments{Username: username})

	if err != nil {
		return nil, err
	}

	return result, err
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
