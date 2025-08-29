package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(".env file not found")
	}
}

func TestPost(t *testing.T) {

	var data Contact
	if err := gofakeit.Struct(&data); err != nil {
		panic(err)
	}

	b, _ := json.Marshal(data)
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
			},
		},
		Body: string(b),
	}

	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}

	var actual Contact
	err := json.Unmarshal([]byte(res.Body), &actual)
	if err != nil {
		t.Error(err)
	}

	t.Log(actual)
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodGet,
			},
		},
	}

	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}

	var actual []Contact
	err := json.Unmarshal([]byte(res.Body), &actual)
	if err != nil {
		t.Error(err)
	}

	for _, v := range actual {
		t.Log(v)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodDelete,
			},
		},
		QueryStringParameters: map[string]string{
			"ids": "01K2Z9913E9DFHHS4AJXEE2S08",
		},
	}

	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}

	t.Log(res.Body)
}
