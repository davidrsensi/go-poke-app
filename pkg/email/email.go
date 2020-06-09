package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "davidrsensi@gmail.com"

	// The HTML body for the email.
	HtmlHeader = "<h2>You were poked by the poke app by one of your friends!</h2>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "You were poked by the poke app by one of your friends!"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Email struct {
	SenderName    string
	Message       string
	ReceiverEmail string
	ReceiverName  string
}

func Send(e Email) {

	var subject string
	var htmlBody string

	if e.SenderName == "Anonymous" {
		subject = "You have been poked"
	} else {
		subject = fmt.Sprintf("%v has poked you!", e.SenderName)
	}

	if e.ReceiverName == "Anonymous" {
		htmlBody = fmt.Sprintf("%v <br> <p>%v</p>", HtmlHeader, e.Message)
	} else {
		htmlBody = fmt.Sprintf("%v <br> <h3>Hey %v,</h3> <img src=\"https://imgur.com/KQPA1PN.png\" width=\"70\" height=\"60\"/> <p>%v</p>", HtmlHeader, e.ReceiverName, e.Message)
	}

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(e.ReceiverEmail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println("Email Sent to address: " + e.ReceiverName)
	fmt.Println(result)
}
