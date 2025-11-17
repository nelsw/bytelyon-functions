package main

import (
	api2 "bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"bytelyon-functions/pkg/service/s3"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db s3.Service

func init() {
	db = s3.New()
}

func handler(req api2.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	if req.Method() == http.MethodOptions {
		return api2.OK()
	}

	v := model.NewBot(req.User())

	switch req.Method() {
	case http.MethodGet:
		return api2.Response(v.FindAll(db))
	case http.MethodPost:
		return api2.Response(v.Create(db, req.Data()))
	case http.MethodDelete:
		return api2.Response(v.Delete(db, req.Param("id")))
	}

	return api2.NotImplemented()
}

func main() {
	lambda.Start(handler)
}
