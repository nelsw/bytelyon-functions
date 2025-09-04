package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func logMap(s string, m map[string]interface{}) {
	log.Info().Any(s, m).Send()
}

func okStr(s string) (events.LambdaFunctionURLResponse, error) {
	return response(http.StatusOK, s)
}

func requestMap(req events.LambdaFunctionURLRequest) map[string]any {
	m := map[string]any{
		"method": req.RequestContext.HTTP.Method,
	}
	if len(req.Headers) > 0 {
		m["headers"] = req.Headers
	}
	if len(req.RequestContext.HTTP.Path) > 0 {
		m["path"] = req.RequestContext.HTTP.Path
	}
	if len(req.QueryStringParameters) > 0 {
		m["query"] = req.QueryStringParameters
	}
	if req.IsBase64Encoded == false && IsJSON(req.Body) {
		var a map[string]any
		_ = json.Unmarshal([]byte(req.Body), &a)
		m["body"] = a
	} else {
		m["body"] = req.Body
	}
	return m
}

func requestMapString(req events.LambdaFunctionURLRequest) string {
	m := requestMap(req)
	b, _ := json.Marshal(&m)
	return string(b)
}

func response(code int, body string) (events.LambdaFunctionURLResponse, error) {

	// log the response so we have full visibility into how the request was handled
	m := map[string]any{"code": code}

	var a map[string]any
	if err := json.Unmarshal([]byte(body), &a); err != nil {
		m["body"] = body
	} else {
		m["body"] = a
	}

	logMap("response", m)

	// return the given ƒ response with a few header values that are required when you QD an API route
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

func BadRequest(err error) (events.LambdaFunctionURLResponse, error) {
	return response(http.StatusBadRequest, err.Error())
}

func LogURLRequest(req events.LambdaFunctionURLRequest) {
	m := map[string]any{
		"headers": req.Headers,
		"method":  req.RequestContext.HTTP.Method,
		"path":    req.RequestContext.HTTP.Path,
		"query":   req.QueryStringParameters,
	}
	if req.IsBase64Encoded == false && IsJSON(req.Body) {
		var a map[string]any
		_ = json.Unmarshal([]byte(req.Body), &a)
		m["body"] = a
	} else {
		m["body"] = req.Body
	}
	logMap("request", requestMap(req))
}

func Marshall(a any) (events.LambdaFunctionURLResponse, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return ServerError(errors.Join(errors.New(fmt.Sprintf("failed to marshal %#v", a)), err))
	}
	return okStr(string(b))
}

func NotImplemented(req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	return response(http.StatusNotImplemented, requestMapString(req))
}

func Err(err error) (events.LambdaFunctionURLResponse, error) {
	if err != nil {
		return BadRequest(err)
	}
	return OK()
}

func OK() (events.LambdaFunctionURLResponse, error) {
	return okStr("")
}

func Response(b []byte, err error) (events.LambdaFunctionURLResponse, error) {

	if err != nil {
		return BadRequest(err)
	}

	var s string
	if b != nil {
		s = string(b)
	}

	return okStr(s)
}

func ServerError(err error) (events.LambdaFunctionURLResponse, error) {
	return response(http.StatusInternalServerError, err.Error())
}

func Unauthorized(err error) (events.LambdaFunctionURLResponse, error) {
	return response(http.StatusUnauthorized, err.Error())
}

func Forbidden() (events.LambdaFunctionURLResponse, error) {
	return response(http.StatusForbidden, "")
}
