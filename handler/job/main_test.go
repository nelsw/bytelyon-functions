package main

import (
	"bytelyon-functions/internal/model/bot"
	"bytelyon-functions/test"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v7"
)

func TestPost(t *testing.T) {
	m := map[string]interface{}{
		"name":     gofakeit.Name(),
		"type":     bot.NewsJobType,
		"keywords": []string{"ford", "bronco"},
		"frequency": map[string]interface{}{
			"unit":  "h",
			"value": 12,
		},
	}

	test.New(t).Post(m).Handle(handler).OK().JSON(m)
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

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(res.Body), &m); err != nil {
		t.Error(err)
	}
	if _, ok := m["items"]; !ok {
		t.Error("items not found")
	}
	size, ok := m["size"]
	if !ok {
		t.Error("size not found")
	}
	n := int(size.(float64))
	if n != len(m["items"].([]interface{})) {
		t.Errorf("got size %d, want %d", n, len(m["items"].([]interface{})))
	}
}
