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

func Handler(req api2.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	v := model.NewProfile(req.User())

	switch req.Method() {
	case http.MethodGet:
		return api2.Response(v, v.Find(db))
	case http.MethodPut:
		if err := v.Hydrate(req.Body); err != nil {
			return api2.BadRequest(err)
		} else if err = v.Validate(); err != nil {
			return api2.BadRequest(err)
		}
		return api2.Response(v, v.Save(db))
	}

	return api2.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
