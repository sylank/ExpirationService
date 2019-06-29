package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.DynamoDBEvent) error {
	fmt.Println(req)

	return nil
}

func main() {
	lambda.Start(handler)
}
