package api

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func Log(v ...interface{}) {
	if v == nil || len(v) == 0 {
		return
	}
	var m = map[string]interface{}{}
	for i := 0; i < len(v); i++ {
		if i%2 == 0 {
			m[v[i].(string)] = nil
		} else {
			m[v[i-1].(string)] = v[i]
		}
	}
	LogMap(m)
}

func LogMap(m map[string]interface{}) {
	b, _ := json.Marshal(m)
	log.Println(string(b))
}

func LogURLRequest(req events.LambdaFunctionURLRequest) {
	LogMap(map[string]interface{}{
		"headers": req.Headers,
		"method":  req.RequestContext.HTTP.Method,
		"query":   req.QueryStringParameters,
		"body":    req.Body,
	})
}

func Response(code int, body string) (events.LambdaFunctionURLResponse, error) {

	// log the response so we have full visibility into how the request was handled
	Log("code", code, "body", body)

	// return the given ƒ response with a few header values that are required when you QD an API route
	// note that we always return nil for error because if we don't, we'll always return a 500
	// as I understand it, it's a worse case scenario akin to self-destruct mode
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
