package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func IsOptions(req events.LambdaFunctionURLRequest) bool {
	return is(req, http.MethodOptions)
}

func IsPost(req events.LambdaFunctionURLRequest) bool {
	return is(req, http.MethodPost)
}

func IsGet(req events.LambdaFunctionURLRequest) bool {
	return is(req, http.MethodGet)
}

func IsPatch(req events.LambdaFunctionURLRequest) bool {
	return is(req, http.MethodPatch)
}

func IsDelete(req events.LambdaFunctionURLRequest) bool {
	return is(req, http.MethodDelete)
}

func is(req events.LambdaFunctionURLRequest, method string) bool {
	return req.RequestContext.HTTP.Method == method
}
