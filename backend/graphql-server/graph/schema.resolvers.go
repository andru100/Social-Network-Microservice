package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"
	"fmt"

	"github.com/andru100/Social-Network-Microservice/backend/graphql-server/model"
)

// SignIn is the resolver for the SignIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.SecurityCheckInput) (*model.Jwtdata, error) {
	panic(fmt.Errorf("not implemented: SignIn - SignIn"))
}

// SignUp is the resolver for the SignUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input model.NewUserDataInput) (*model.Jwtdata, error) {
	panic(fmt.Errorf("not implemented: SignUp - SignUp"))
}

// LikeComment is the resolver for the LikeComment field.
func (r *mutationResolver) LikeComment(ctx context.Context, input model.SendLikeInput) (*model.MongoFields, error) {
	panic(fmt.Errorf("not implemented: LikeComment - LikeComment"))
}

// ReplyComment is the resolver for the ReplyComment field.
func (r *mutationResolver) ReplyComment(ctx context.Context, input model.ReplyCommentInput) (*model.MongoFields, error) {
	panic(fmt.Errorf("not implemented: ReplyComment - ReplyComment"))
}

// NewComment is the resolver for the NewComment field.
func (r *mutationResolver) NewComment(ctx context.Context, input model.SendCmtInput) (*model.MongoFields, error) {
	panic(fmt.Errorf("not implemented: NewComment - NewComment"))
}

// UpdateBio is the resolver for the UpdateBio field.
func (r *mutationResolver) UpdateBio(ctx context.Context, input model.UpdateBioInput) (*model.MongoFields, error) {
	panic(fmt.Errorf("not implemented: UpdateBio - UpdateBio"))
}

// RequestOtp is the resolver for the RequestOTP field.
func (r *mutationResolver) RequestOtp(ctx context.Context, input model.RequestOtpInput) (*model.Confirmation, error) {
	panic(fmt.Errorf("not implemented: RequestOtp - RequestOTP"))
}

// SecureUpdate is the resolver for the SecureUpdate field.
func (r *mutationResolver) SecureUpdate(ctx context.Context, input model.SecurityCheckInput) (*model.Jwtdata, error) {
	panic(fmt.Errorf("not implemented: SecureUpdate - SecureUpdate"))
}

// Chkauth is the resolver for the Chkauth field.
func (r *queryResolver) Chkauth(ctx context.Context, input model.JwtdataInput) (*model.Authd, error) {
	panic(fmt.Errorf("not implemented: Chkauth - Chkauth"))
}

// GetAllComments is the resolver for the GetAllComments field.
func (r *queryResolver) GetAllComments(ctx context.Context, input string) (*model.MongoFields, error) {
	panic(fmt.Errorf("not implemented: GetAllComments - GetAllComments"))
}

// GetUserComments is the resolver for the GetUserComments field.
func (r *queryResolver) GetUserComments(ctx context.Context, input string) (*model.MongoFields, error) {
	panic(fmt.Errorf("not implemented: GetUserComments - GetUserComments"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
