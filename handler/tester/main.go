package main

import (
	"bytelyon-functions/pkg/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	return api.Marshall(req)
}

func main() {
	lambda.Start(Handler)
}
