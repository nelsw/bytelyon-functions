package main

import (
	"bytelyon-functions/pkg/api"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	if api.IsOptions(req) {
		return api.OK()
	} else if api.IsPost(req) {

	}
	return api.NotImplemented("Method not implemented: " + req.RequestContext.HTTP.Method)
}

func handleLogin(ctx context.Context, token string) (events.LambdaFunctionURLResponse, error) {

	//c, err := model.NewCredentials(token)
	//if err != nil {
	//	return api.Response(http.StatusBadRequest, err.Error())
	//}
	//
	//var email model.Email
	//if err = entity.New(ctx).Value(&email).ID(c.Email).Find(); err != nil {
	//	return api.Response(http.StatusBadRequest, "email not found")
	//}

	return api.OK()

	//return api.Response(http.StatusOK, auth.NewToken(map[string]interface{}{
	//	"email":          u.Email,
	//	"email_verified": email.Verified,
	//	"name":           up.Name,
	//	"image":          up.Image,
	//}))
}

func handleSignup(ctx context.Context, token string) (events.LambdaFunctionURLResponse, error) {

	//c, err := model.NewCredentials(token)
	//if err != nil {
	//	return api.Response(http.StatusBadRequest, err.Error())
	//}
	//
	//u := c.NewUser()
	//if err = entity.New(ctx).Value(u).Save(); err != nil {
	//	return api.Response(http.StatusInternalServerError, err.Error())
	//} else if err = entity.New(ctx).Value(c.NewUserProfile(u.ID)).Save(); err != nil {
	//	return api.Response(http.StatusInternalServerError, err.Error())
	//} else if err = entity.New(ctx).Value(c.NewPassword(u.ID)).Save(); err != nil {
	//	return api.Response(http.StatusInternalServerError, err.Error())
	//}
	//
	//e := c.NewEmail(u.ID)
	//if err = entity.New(ctx).Value(e).Save(); err != nil {
	//	return api.Response(http.StatusInternalServerError, err.Error())
	//} else if err = ses.New(ctx).VerifyEmail(ctx, c.Email, e.Token); err != nil {
	//	return api.Response(http.StatusInternalServerError, err.Error())
	//}

	return api.OK()
}

func main() {
	lambda.Start(handler)
}
