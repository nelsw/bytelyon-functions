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
		"name":     gofakeit.Name(),
		"type":     bot.GoogleNews,
		"keywords": []string{"ford", "bronco"},
		"frequency": map[string]interface{}{
			"unit":  "h",
			"value": 12,
		},
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
