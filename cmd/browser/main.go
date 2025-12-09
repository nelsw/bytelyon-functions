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

	v := model.NewBrowser(req.User())

	switch req.Method() {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		return api.Response(v.Save([]byte(req.Body)))
	case http.MethodGet:
		return api.Response(v.Find())
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
