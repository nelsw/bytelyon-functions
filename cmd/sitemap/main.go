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

	switch req.Method() {
	case http.MethodGet:
		return api.Response(req.User().Sitemaps())
	case http.MethodPost:
		return api.Response(model.NewSitemap(req.User()).Create([]byte(req.Body)))
	case http.MethodDelete:
		return api.Response(model.NewSitemap(req.User(), req.Param("domain"), req.Param("id")).Delete())
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
