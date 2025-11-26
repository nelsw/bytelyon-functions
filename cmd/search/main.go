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

	v, err := model.NewSitemap(req.User(), req.Param("url"))
	if err != nil {
		return api.BadRequest(err)
	}

	switch req.Method() {
	case http.MethodPost:
		return api.Response(v.Create())
	case http.MethodGet:
		api.NotImplemented()
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
