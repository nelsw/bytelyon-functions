package main

import (
	"bytelyon-functions/pkg/api"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	b, _ := json.MarshalIndent(req, "", "\t")

	return api.Response(http.StatusOK, string(b))
}

func main() {
	lambda.Start(handler)
}
