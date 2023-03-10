package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Uploaditem(bucket *string, filename *string, filebytes *[]byte) (*string, error) { // upload file to s3 with the bucket name and file adress passed to it

	tmpfile, err := ioutil.TempFile("", "example") // create temp file using naming convention.. it'll ad random stuff
	
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(*filebytes); err != nil { //write file from []bytes given by io readall
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(tmpfile.Name()) // open file using the temp dir and temp name created.
	if err != nil {
		fmt.Println("Unable to open file %q, %v", err)
	}

	defer file.Close() // clean up

	uploader := s3manager.NewUploader(Sess)

	result, err := uploader.Upload(&s3manager.UploadInput{ // upload file
		Bucket: aws.String(*bucket + Uniqueadr),
		Key:    aws.String(*filename),
		Body:   file,
	})

	if err != nil {
		return nil, errors.New("Unable to upload file")
	}

	fmt.Printf("Successfully uploaded %q to %q\n", *filename, *bucket+Uniqueadr)
	
	return &result.Location, err

}
