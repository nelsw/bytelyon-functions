package main

import (
	"bytelyon-functions/internal/model/bot"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v7"
)

func TestPost(t *testing.T) {
	b, _ := json.Marshal(map[string]interface{}{
		"name":      gofakeit.Name(),
		"type":      bot.GoogleNews,
		"roots":     []string{},
		"keywords":  []string{"ford", "bronco"},
		"frequency": 60,
	})
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
			},
		},
		QueryStringParameters: map[string]string{},
		Body:                  string(b),
	}

	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("got status code %d, want %d", res.StatusCode, http.StatusOK)
	}

	var v bot.Job
	err := json.Unmarshal([]byte(res.Body), &v)
	if err != nil {
		t.Error(err)
	}
	b, _ = json.MarshalIndent(&v, "", "\t")
	t.Log(string(b))
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "GET",
			},
		},
		QueryStringParameters: map[string]string{
			"type": "job",
		},
	}
	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("got status code %d, want %d", res.StatusCode, http.StatusOK)
	}

	var vv []bot.Job
	err := json.Unmarshal([]byte(res.Body), &vv)
	if err != nil {
		t.Error(err)
	}
	b, _ := json.MarshalIndent(&vv, "", "\t")
	t.Log(string(b))
}
