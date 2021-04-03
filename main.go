package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

var sess *session.Session

// Basic information for the Amazon Elasticsearch Service domain
var endpoint = "https://search-cupcake-domain-001-bdj74ottahj7ttzw3szp4e6tea.ap-south-1.es.amazonaws.com"
var region = "ap-south-1" // e.g. us-east-1
var service = "es"
var signer *v4.Signer
var client = &http.Client{}

// this runs only the very first time
// this lambda function is created
// will be used for initial setup
func init() {
	sess = configureAWS()
	signer = v4.NewSigner(sess.Config.Credentials)
}

func configureAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials("AKIATNLFYBIN5NZE2OXK", "WclXT3YX5rqQksVNznG5At7IL+haRkak5vD3eMri", ""),
	})

	if err != nil {
		log.Fatal(err)
	}

	// adding logging handler to session
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		// Log every request made and its payload
		fmt.Printf("Request: %s/%s, Payload: %s",
			r.ClientInfo.ServiceName, r.Operation, r.Params)
	})

	// fmt.Println(sess)

	return sess
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
		fmt.Printf("Change occured: %s", record.Change)

		// Print new values for attributes of type String
		for name, value := range record.Change.NewImage {
			if value.DataType() == events.DataTypeString {
				fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
			}
		}
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, body)
	if err != nil {
		fmt.Print(err)
	}

}

func main() {
	lambda.Start(handleRequest)
}
