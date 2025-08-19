package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/internal/service/auth"
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/service/ses"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	switch req.RequestContext.HTTP.Method {
	case http.MethodOptions:
		return api.Response(http.StatusOK, "")
	case http.MethodPost:
		if req.QueryStringParameters["forgot"] == "true" {
			return api.Response(http.StatusNotImplemented, "Forgot password not implemented")
		}
		if req.QueryStringParameters["login"] == "true" {
			return handleLogin(ctx, req.Headers["authorization"])
		}
		if req.QueryStringParameters["signup"] == "true" {
			return handleLogin(ctx, req.Headers["authorization"])
		}
		fallthrough
	default:
		return api.Response(http.StatusNotImplemented, "Method not implemented: "+req.RequestContext.HTTP.Method)
	}
}

func handleLogin(ctx context.Context, token string) (events.LambdaFunctionURLResponse, error) {

	c, err := model.NewCredentials(token)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	var email model.UserEmail
	if err = entity.New(ctx).Value(&email).ID(c.Email).Find(); err != nil {
		return api.Response(http.StatusBadRequest, "email not found")
	}

	var password model.UserPassword
	if err = entity.New(ctx).Value(&password).ID(email.UserID).Find(); err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	} else if err = password.Validate(c.Password); err != nil {
		return api.Response(http.StatusUnauthorized, "incorrect password")
	}

	return api.Response(http.StatusOK, auth.NewToken(email.UserID))
}

func handleSignup(ctx context.Context, token string) (events.LambdaFunctionURLResponse, error) {

	c, err := model.NewCredentials(token)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	if err = entity.New(ctx).Value(c.NewUser()).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	e := c.NewEmail()
	if err = entity.New(ctx).Value(c.NewEmail()).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	if err = entity.New(ctx).Value(c.NewPassword()).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	if err = ses.New(ctx).VerifyEmail(ctx, c.Email, e.Token); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	return api.Response(http.StatusOK, "")
}

func main() {
	lambda.Start(handler)
}
