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

	v := model.NewNews(req.User())

	switch req.Method() {
	case http.MethodDelete:
		return api.Response(nil, v.Delete())
	case http.MethodPost:
		return api.Response(v.Create([]byte(req.Body)))
	case http.MethodGet:
		return api.Response(v.FindAll())
	case http.MethodPatch:
		v.Work()
		return api.OK()
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
