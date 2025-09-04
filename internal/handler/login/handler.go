package login

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	token := req.Headers["authorization"]
	if strings.HasPrefix(token, "Basic ") {
		token = strings.TrimPrefix(token, "Basic ")
	}

	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return app.BadRequest(err)
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return app.BadRequest(errors.New("invalid token"))
	}

	db := s3.NewWithContext(ctx)

	email := model.Email{ID: parts[0]}
	if b, err = db.Get(email.Path()); err != nil {
		return app.BadRequest(err)
	}
	app.MustUnmarshal(b, &email)

	pass := model.Password{User: email.User(), Text: parts[1]}
	if b, err = db.Get(pass.Path()); err != nil {
		return app.BadRequest(err)
	}
	app.MustUnmarshal(b, &pass)

	if err = pass.Compare(); err != nil {
		return app.Unauthorized(errors.Join(err, errors.New("invalid password")))
	}

	return app.Response(model.CreateJWT(ctx, email.User()))
}
