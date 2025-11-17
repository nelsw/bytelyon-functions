package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Response(a any, err error) (events.APIGatewayV2HTTPResponse, error) {

	if err != nil {
		return Error(http.StatusInternalServerError, err)
	}

	return OK(a)
}

func OK(aa ...any) (events.APIGatewayV2HTTPResponse, error) {

	var a any
	if len(aa) > 0 {
		a = aa[0]
	}

	if a == nil {
		return response(http.StatusOK, "")
	}

	b, err := json.Marshal(a)
	if err != nil {
		return Error(http.StatusInternalServerError, err)
	}

	return response(http.StatusOK, string(b))
}

func NotImplemented() (events.APIGatewayV2HTTPResponse, error) {
	return Error(http.StatusNotImplemented, nil)
}

func BadRequest(err error) (events.APIGatewayV2HTTPResponse, error) {
	return Error(http.StatusBadRequest, err)
}

func Error(code int, err error) (events.APIGatewayV2HTTPResponse, error) {

	if err == nil {
		err = errors.New(http.StatusText(code))
	}

	return response(code, err.Error())
}

func response(code int, body string) (events.APIGatewayV2HTTPResponse, error) {

	var lvl zerolog.Level
	if code < 300 {
		lvl = zerolog.InfoLevel
	} else if code < 500 {
		lvl = zerolog.WarnLevel
	} else {
		lvl = zerolog.ErrorLevel
	}

	log.WithLevel(lvl).
		Int("code", code).
		Str("body", body).
		Msg("API Response")

	return events.APIGatewayV2HTTPResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "authorization, content-type,",
			"Access-Control-Allow-Methods": "*",
		},
		StatusCode: code,
		Body:       body,
	}, nil
}
