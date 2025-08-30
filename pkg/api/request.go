package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func IsOptions(req events.LambdaFunctionURLRequest) bool {
	return req.RequestContext.HTTP.Method == http.MethodOptions
}

func IsPost(req events.LambdaFunctionURLRequest) bool {
	return req.RequestContext.HTTP.Method == http.MethodPost
}
