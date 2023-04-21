package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func Connectedmngo(err error, err1 error) { // prints connected if all error checks passed
	if err != nil || err1 != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}
}

func Connectedaws(err2 error) { // prints connected to aws if all error checks passed
	if err2 != nil {
		log.Fatal(err2)
	} else {
		fmt.Println("Connected to MongoDB!")
	}
}

func Listbuckets() {
	// Create S3 service client
	svc := s3.New(Sess)

	// list all buckets
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets { // loop through results print name and date created
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func Createbucket(bucketname string) { // creates a s3 bucket with the name passed to it
	// Create S3 service client
	svc := s3.New(Sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketname + Uniqueadr), // make bucket name unique
	})
	if err != nil {
		fmt.Println("Unable to create bucket %q, %v", bucketname+Uniqueadr, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucketname+Uniqueadr)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketname + Uniqueadr),
	})

	if err != nil { // check bucket is created
		fmt.Println("Error occurred while waiting for bucket to be created, %v", bucketname+Uniqueadr)
	} else {
		fmt.Printf("Bucket %q successfully created\n", bucketname+Uniqueadr)
	}

	Publicbucket(bucketname) // make bucket public read only
}

func Publicbucket(bucket string) { // make bucket public and read only
	// Create S3 service client
	svc := s3.New(Sess)

	readOnlyAnonUserPolicy := map[string]interface{}{ // add policy to map so can be sent
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":       "AddPerm",
				"Effect":    "Allow",
				"Principal": "*",
				"Action": []string{
					"s3:GetObject",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", bucket+Uniqueadr),
				},
			},
		},
	}

	policy, err := json.Marshal(readOnlyAnonUserPolicy)

	if err != nil {
		fmt.Println("Unable to marshal json %v", err)
	}

	_, err = svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket + Uniqueadr),
		Policy: aws.String(string(policy)),
	})

	if err != nil {
		fmt.Println("Unable to update bucket policy %v", err)
	} else {
		fmt.Printf("Successfully set bucket %q's policy\n", bucket+Uniqueadr)
	}

}

func Listitems(bucket string) {
	// Create S3 service client
	svc := s3.New(Sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		fmt.Println("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents { // loop through bucket contents and use attributes to print info
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

}

func GetTempLoc(filename string) string { // get the temp location of files sent in post requests to be used for aws upload without having to store data locally
	return strings.TrimRight(os.TempDir(), "/") + "/" + filename
}

func Getfiletype(filename string) string { // takes filname and returns filetype

	runed := []rune(filename)
	var result [][]rune

	for i := len(runed) - 1; i >= 0; i-- { // looks for . and cuts the file type to new array
		if runed[i] == 46 {
			result = append(result, runed[i:])
			break
		}
	}

	return (string(result[0]))
}

func ExitErrorf(msg string, args ...interface{}) { // hadles aws errors
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func CORSMiddleware() gin.HandlerFunc { // cors func to allow body and acces from anywhere
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
