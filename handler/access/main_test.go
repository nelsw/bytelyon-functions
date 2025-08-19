package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestLogin(t *testing.T) {

	data := []byte("kowalski7012@gmail.com:Farts1234!")
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
			},
		},
		Headers: map[string]string{
			"authorization": "Basic " + base64.StdEncoding.EncodeToString(data),
		},
		QueryStringParameters: map[string]string{
			"login": "true",
		},
	}

	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}

}
