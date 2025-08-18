package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/service/s3"
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/ulid/v2"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.Log("context", ctx, "request", req)

	switch req.RequestContext.HTTP.Method {
	case http.MethodPost:
		return handlePost(ctx, req.Body)
	case http.MethodOptions:
		return api.Response(http.StatusOK, "")
	default:
		return api.Response(http.StatusNotImplemented, "Method not implemented: "+req.RequestContext.HTTP.Method)
	}
}

func handlePost(ctx context.Context, body string) (events.LambdaFunctionURLResponse, error) {

	// check that the given body actually contains data
	if len(body) == 0 {
		return api.Response(http.StatusBadRequest, "")
	}

	// use the bucket defined in our env vars
	bucket := os.Getenv("BUCKET")

	// use a key that's guaranteed to be unique but also sortable
	// include a json file extension so we can read what we save
	key := ulid.Make().String() + ".json"

	// convert the given request body string to bytes
	data := []byte(body)

	// try to put the data and return a 500 with the error message if it fails
	if err := s3.NewClient(ctx).Put(ctx, bucket, key, data); err != nil {
		return api.Response(http.StatusInternalServerError, "While putting data: "+err.Error())
	}

	return api.Response(http.StatusOK, "")
}

func main() {
	lambda.Start(handler)
}
