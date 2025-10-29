package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

const msg = "API Response"

func Response(a any, err error) (events.APIGatewayV2HTTPResponse, error) {

	if err != nil {
		return Error(http.StatusInternalServerError, err)
	}

	log.Info().
		Int("code", http.StatusOK).
		Any("body", a).
		Msg(msg)

	if a == nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
		}, nil
	}

	var b []byte
	if b, err = json.Marshal(a); err != nil {
		return Error(http.StatusInternalServerError, err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(b),
	}, nil
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

	log.Err(err).
		Int("code", code).
		Msg(msg)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Body:       err.Error(),
	}, nil
}
