package main

import (
	"bytelyon-functions/internal/api"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/internal/service/s3"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db s3.Service

func init() {
	db = s3.New()
}

func handler(req api.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	if req.Method() == http.MethodOptions {
		return api.OK()
	}

	v := model.NewBot(req.User())

	switch req.Method() {
	case http.MethodGet:
		return api.Response(v.FindAll(db))
	case http.MethodPost:
		return api.Response(v.Create(db, req.Data()))
	case http.MethodDelete:
		return api.Response(v.Delete(db, req.Param("id")))
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(handler)
}
