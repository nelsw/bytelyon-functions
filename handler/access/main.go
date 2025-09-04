package main

import (
	"bytelyon-functions/internal/auth"
	"bytelyon-functions/pkg/api"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	if api.IsOptions(req) {
		return api.OK()
	}

	if api.IsPost(req) {
		switch req.RawPath {
		case "/login":
			return api.Response(auth.Login(ctx, req.Headers["authorization"]))
		}
	}

	return api.NotImplemented(req)
}

func main() {
	lambda.Start(handler)
}
