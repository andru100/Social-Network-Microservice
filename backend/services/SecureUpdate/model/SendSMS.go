package model

import (
	"fmt"

	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/andru100/Social-Network-Microservices/backend/services/SecureUpdate/utils"
)

func SendSMS(mobile *string, otp *string) (bool, error) {

	// Create SNS service client
	svc := sns.New(utils.Sess)

	// Build the request input
	params := &sns.PublishInput{
		Message: aws.String("Your otp is " + *otp),
		PhoneNumber: aws.String(*mobile),
	}

	// Send the request
	resp, err := svc.Publish(params)

	if err != nil {
		fmt.Println(resp, err)
	}

	return true, err
}

