package model

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sms"
	"github.com/aws/aws-sdk-go/aws"
)

func SendSMS(mobile string, otp string) (string, error) {

	// Create SNS service client
	svc := sms.New(Sess)

	// Build the request input
	params := &sms.SendSMSInput{
		Message: aws.String("Your otp is " + otp),
		PhoneNumber: aws.String(mobile),
	}

	// Send the request
	resp, err := svc.SendSMS(params)

	// Log the response
	fmt.Println(resp)

	return resp, err
}

