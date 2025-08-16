package main

import (
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/ulid/v2"
)

// bucket is the existing s3 used to store contact form data
const bucket = "bytelyon-contact"

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	// log the given arguments in a format that CloudWatch
	// will make readable without breaking the bank
	b, _ := json.Marshal(map[string]interface{}{
		"ctx": ctx,
		"req": req,
	})
	log.Println(string(b))

	// return 200 for preflight requests
	if req.RequestContext.HTTP.Method == http.MethodOptions {
		return response(http.StatusOK, "")
	}

	// check that the given body actually contains data
	if len(req.Body) == 0 {
		return response(http.StatusBadRequest, "")
	}

	// we've got some bytes, lets save em to s3
	s3Client := s3.NewClient(ctx)

	// use a key that's guaranteed to be unique but also sortable
	// include a json file extension so we can read what we save
	key := ulid.Make().String() + ".json"

	// convert the given request body string to bytes
	data := []byte(req.Body)

	// try to put the data and return a 500 with the error message if it fails
	if err := s3Client.Put(ctx, bucket, key, data); err != nil {
		return response(http.StatusInternalServerError, "While putting data: "+err.Error())
	}

	// gravy train with biscuit wheels ... return success
	return response(http.StatusOK, "")
}

// response is a helper ƒ for returning a events.LambdaFunctionURLResponse with necessary headers
func response(code int, body string) (events.LambdaFunctionURLResponse, error) {

	// log the response so we have full visibility into how the request was handled
	b, _ := json.Marshal(map[string]interface{}{
		"code": code,
		"body": body,
	})
	log.Println(string(b))

	// return the given ƒ response with a few header values that are required when you QD an API route
	// note that we always return nil for error because if we don't, we'll always return a 500
	// as I understand it, it's a worse case scenario akin to self destruct mode
	return events.LambdaFunctionURLResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Authorization,Content-Type",
			"Access-Control-Allow-Methods": "*",
		},
		StatusCode: code,
		Body:       body,
	}, nil
}

func main() {
	lambda.Start(handler)
}
