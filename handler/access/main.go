package main

import (
	"bytelyon-functions/internal/model"
	"bytelyon-functions/internal/util"
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	if api.IsOptions(req) {
		return api.OK()
	}

	if api.IsPost(req) {
		switch req.RawPath {
		case "/login":
			return api.Response(handleLogin(ctx, req.Headers["authorization"]))
		}
	}

	return api.NotImplemented(req)
}

func handleLogin(ctx context.Context, token string) ([]byte, error) {
	if strings.HasPrefix(token, "Basic ") {
		token = strings.TrimPrefix(token, "Basic ")
	}

	b, _ := base64.StdEncoding.DecodeString(token)
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}

	db := s3.NewWithContext(ctx)

	email := model.Email{ID: parts[0]}
	if out, err := db.Get(ctx, "bytelyon", email.Key()); err != nil {
		return nil, err
	} else {
		util.MustUnmarshal(out, &email)
	}

	pass := model.Password{UserID: email.UserID, Text: parts[1]}
	if out, err := db.Get(ctx, "bytelyon", pass.Key()); err != nil {
		return nil, err
	} else {
		util.MustUnmarshal(out, &pass)
		if err = pass.Compare(); err != nil {
			return nil, err
		}
	}

	return model.CreateJWT(ctx, model.User{ID: email.UserID})
}

func main() {
	lambda.Start(handler)
}
