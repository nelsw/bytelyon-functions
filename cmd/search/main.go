package main

import (
	"bytelyon-functions/pkg/api"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req api.Request) (events.APIGatewayV2HTTPResponse, error) {
	req.Log()

	switch req.Method() {
	case http.MethodGet:
		return api.Response(req.User().Searches())
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
