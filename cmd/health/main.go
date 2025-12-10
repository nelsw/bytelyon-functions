package main

import (
	"bytelyon-functions/pkg/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req api.Request) (events.APIGatewayV2HTTPResponse, error) {
	req.Log()
	return api.OK(req)
}

func main() {
	lambda.Start(Handler)
}
