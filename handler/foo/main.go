package main

import (
	"bytelyon-functions/pkg/api"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	api.LogURLRequest(req)
	return api.OK(os.Getenv("FOO"))
}

func main() {
	lambda.Start(handler)
}
