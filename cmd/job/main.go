package main

import (
	"bytelyon-functions/internal/api"
	"bytelyon-functions/internal/model"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req api.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	if req.RequestContext.HTTP.Method == http.MethodOptions {
		return api.OK()
	}

	v, err := model.NewJob(req.User(), req.Param("id"), req.Data())
	if err != nil {
		return api.BadRequest(err)
	}

	switch req.Method() {
	case http.MethodDelete:
		return api.Response(v.Delete())
	case http.MethodPost:
		return api.Response(v.Create())
	case http.MethodPut:
		return api.Response(v.Update())
	case http.MethodGet:
		return api.Response(v.FindAll())
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
