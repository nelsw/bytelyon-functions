package main

import (
	api2 "bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req api2.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	if req.Method() == http.MethodOptions {
		return api2.OK()
	}

	v, err := model.NewSitemap(req.User(), req.Param("url"))
	if err != nil {
		return api2.BadRequest(err)
	}

	switch req.Method() {
	case http.MethodGet:
		return api2.Response(v.FindAll())
	case http.MethodPost:
		return api2.Response(v.Create())
	case http.MethodDelete:
		return api2.Response(v.Delete())
	}

	return api2.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
