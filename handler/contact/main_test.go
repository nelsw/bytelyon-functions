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

func TestHandler(t *testing.T) {
	var data struct {
		Name    string `json:"name" fake:"{name}"`
		Email   string `json:"email" fake:"{email}"`
		Message string `json:"message" fake:"{sentence}"`
	}
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
}
