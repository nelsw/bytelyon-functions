package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req api.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	v := model.NewProfile(req.User())

	switch req.Method() {
	case http.MethodGet:
		return api.Response(v.Find())
	case http.MethodPut:
		return api.Response(v.Create([]byte(req.Body)))
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
