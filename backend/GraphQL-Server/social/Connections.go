package social

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect db globally so all funcs can use client
var ClientOptions = options.Client().ApplyURI("mongodb+srv://andru:1q1q1q@cluster0.tccti.mongodb.net/cluster0?retryWrites=true&w=majority") // Set client options

var Client, err = mongo.Connect(context.TODO(), ClientOptions) // Connect to MongoDB

var Err1 = Client.Ping(context.TODO(), nil) // Check the connection

var Sess, Err2 = session.NewSession(&aws.Config{ //start a aws session by setting the region
	Region: aws.String("eu-west-2")},
)

var Uniqueadr = "ajh46unique" // used to make s3 file upload filenames unique
