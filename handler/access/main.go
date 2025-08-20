package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model/auth"
	"bytelyon-functions/internal/model/credentials"
	"bytelyon-functions/internal/model/user"
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
		} else if req.QueryStringParameters["login"] == "true" {
			return handleLogin(ctx, req.Headers["authorization"])
		} else if req.QueryStringParameters["signup"] == "true" {
			return api.Response(http.StatusNotImplemented, "Signup not implemented")
		}
		fallthrough
	default:
		return api.Response(http.StatusNotImplemented, "Method not implemented: "+req.RequestContext.HTTP.Method)
	}
}

func handleLogin(ctx context.Context, token string) (events.LambdaFunctionURLResponse, error) {

	c, err := credentials.NewCredentials(token)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	var email user.UserEmail
	if err = entity.New(ctx).Value(&email).ID(c.Email).Find(); err != nil {
		return api.Response(http.StatusBadRequest, "email not found")
	}

	var password user.UserPassword
	if err = entity.New(ctx).Value(&password).ID(email.UserID).Find(); err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	} else if err = password.Validate(c.Password); err != nil {
		return api.Response(http.StatusUnauthorized, "incorrect password")
	}

	var u user.User
	if err = entity.New(ctx).Value(&u).ID(email.UserID).Find(); err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	var up user.UserProfile
	_ = entity.New(ctx).Value(&up).ID(email.UserID).Find()

	return api.Response(http.StatusOK, auth.NewToken(map[string]interface{}{
		"email":          u.Email,
		"email_verified": email.Verified,
		"name":           up.Name,
		"image":          up.Image,
	}))
}

func handleSignup(ctx context.Context, token string) (events.LambdaFunctionURLResponse, error) {

	c, err := credentials.NewCredentials(token)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	u := c.NewUser()
	if err = entity.New(ctx).Value(u).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	} else if err = entity.New(ctx).Value(c.NewUserProfile(u.ID)).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	} else if err = entity.New(ctx).Value(c.NewPassword(u.ID)).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	e := c.NewEmail(u.ID)
	if err = entity.New(ctx).Value(e).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	} else if err = ses.New(ctx).VerifyEmail(ctx, c.Email, e.Token); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	return api.Response(http.StatusOK, "")
}

func main() {
	lambda.Start(handler)
}
