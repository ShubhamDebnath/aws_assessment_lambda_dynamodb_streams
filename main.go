package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
var domain = "https://search-cupcake-domain-001-bdj74ottahj7ttzw3szp4e6tea.ap-south-1.es.amazonaws.com"
var index = "cupcake-index-002"
var endpoint = domain + "/" + index + "/" + "_doc" + "/"
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

	maps := make([]map[string]string, 0)

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
		attrValMap := make(map[string]string)

		// Print new values for attributes of type String
		for name, value := range record.Change.NewImage {
			if value.DataType() == events.DataTypeString {
				fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
				attrValMap[name] = value.String()

			} else if value.DataType() == events.DataTypeNumber {
				fmt.Printf("Numerical Attribute name: %s, value: %s\n", name, value.Number())
				attrValMap[name] = value.Number()
			}
		}
		attrValMap["evtType"] = record.EventName
		maps = append(maps, attrValMap)
	}

	for _, map1 := range maps {
		jsonBody, err := json.Marshal(map1)
		body := strings.NewReader(string(jsonBody))

		fmt.Println(string(jsonBody))

		if err != nil {
			fmt.Println(err)
		}

		req, err := http.NewRequest(http.MethodPost, endpoint, body)
		if err != nil {
			fmt.Print(err)
		}

		req.Header.Add("Content-Type", "application/json")

		signer.Sign(req, body, service, region, time.Now())
		resp, err := client.Do(req)

		if err != nil {
			fmt.Print(err)
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		fmt.Println(string(bodyBytes))
	}
}

func main() {
	lambda.Start(handleRequest)
}
