package app

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func IsDelete(req events.LambdaFunctionURLRequest) bool  { return is(req, http.MethodDelete) }
func IsGet(req events.LambdaFunctionURLRequest) bool     { return is(req, http.MethodGet) }
func IsOptions(req events.LambdaFunctionURLRequest) bool { return is(req, http.MethodOptions) }
func IsPatch(req events.LambdaFunctionURLRequest) bool   { return is(req, http.MethodPatch) }
func IsPost(req events.LambdaFunctionURLRequest) bool    { return is(req, http.MethodPost) }
func IsPut(req events.LambdaFunctionURLRequest) bool     { return is(req, http.MethodPut) }
func is(req events.LambdaFunctionURLRequest, method string) bool {
	return req.RequestContext.HTTP.Method == method
}
