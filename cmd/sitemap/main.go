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

	if req.Method() == http.MethodOptions {
		return api.OK()
	}

	v, err := model.NewSitemap(req.User(), req.Param("url"))
	if err != nil {
		return api.BadRequest(err)
	}

	switch req.Method() {
	case http.MethodGet:
		return api.Response(v.FindAll())
	case http.MethodPost:
		return api.Response(v.Create())
	case http.MethodDelete:
		return api.Response(v.Delete())
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
