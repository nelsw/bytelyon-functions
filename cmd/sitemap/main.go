package main

import (
	"bytelyon-functions/internal/api"
	"bytelyon-functions/internal/model"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	v, err := model.NewSitemap(req)
	if err != nil {
		return api.BadRequest(err)
	}

	switch req.RequestContext.HTTP.Method {
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
