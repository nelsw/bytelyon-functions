package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(r api.Request) (events.APIGatewayV2HTTPResponse, error) {
	r.Log()

	v := model.NewArticle(r.User(), r.Param("news"), r.Param("id"))

	switch r.Method() {
	case http.MethodDelete:
		return api.Response(nil, v.Delete())
	case http.MethodGet:
		if v.ID.IsZero() {
			return api.Response(v.FindAll())
		}
		return api.Response(v, v.Find())
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
